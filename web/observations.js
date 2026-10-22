// ==================================================
// 1. HTML要素の取得
// ==================================================

const nextPassCity = document.getElementById("next-pass-city");
const nextPassCountdown = document.getElementById("next-pass-countdown");
const nextPassStart = document.getElementById("next-pass-start");
const nextPassDirection = document.getElementById("next-pass-direction");
const nextPassElevation = document.getElementById("next-pass-elevation");
const nextPassCloud = document.getElementById("next-pass-cloud");
const nextPassVisibility = document.getElementById("next-pass-visibility");
const nextPassMessage = document.getElementById("next-pass-message");
// 右側
const passesList = document.getElementById("passes-list");
const passesSummary = document.getElementById("passes-summary");
const passesPageInfo = document.getElementById("passes-page-info");
const passesPrev = document.getElementById("passes-prev");
const passesNext = document.getElementById("passes-next");
// sort関連
const visibilityFilter = document.getElementById("visibility-filter");
const elevationFilter = document.getElementById("elevation-filter");
const passesSort = document.getElementById("passes-sort");

// 現在のページ
let currentPage = 1;

// 古いリクエストの結果で画面が上書きされないようにする
let passesRequestId = 0;

// カウントダウンのタイマーを管理
let countdownTimer = null;
// ==================================================
// 2. 通過情報の取得（現在はモック）
// ==================================================
// GET /api/iss/passesを呼ぶ cityIDにはtokyoとかが入る
async function fetchISSPasses(cityId, options = {}) {
    const {
        visibility = "",
        minElevation = 0,
        sort = "start",
        page = 1,
        limit = 3,
    } = options;

    // --------------------------------------------------
    // TODO: 本番用：API完成後にコメントを外す
    // --------------------------------------------------

    // // 都市IDをクエリパラメータにする
    // const params = new URLSearchParams({
    //     city: cityId,
    //     minElevation: String(minElevation),
    //     sort: sort,
    //     page: String(page),
    //     limit: String(limit),
    // });

    // if (visibility !== "") {
    //     params.set("visibility", visibility);
    // }

    // // 自分たちのGo APIへリクエストを送る
    // const response = await fetch(`${API_BASE_URL}/api/iss/passes?${params}`);

    // // HTTPエラーを確認する
    // if (!response.ok) {
    //     throw new Error(
    //         `通過情報の取得に失敗しました: HTTP ${response.status}`,
    //     );
    // }
    // // レスポンス本文のJSONを読み取る
    // const data = await response.json();

    // return data;

    // --------------------------------------------------
    // TODO: モック用：API接続時は、この部分を削除またはコメントアウト
    // --------------------------------------------------

    const now = Date.now();

    // 各通過の仮の条件
    const mockConditions = [
        {
            hoursLater: 2,
            elevation: 62,
            cloudCover: 20,
            magnitude: -3,
            visibility: "good",
            reasons: ["雲量が少ない", "最大仰角が高い"],
            message: "空の高い位置を通り、観測しやすい条件です。",
        },
        {
            hoursLater: 26,
            elevation: 35,
            cloudCover: 45,
            magnitude: -1.5,
            visibility: "fair",
            reasons: ["雲量が中程度"],
            message: "雲の切れ間から見える可能性があります。",
        },
        {
            hoursLater: 50,
            elevation: 15,
            cloudCover: 80,
            magnitude: -1,
            visibility: "poor",
            reasons: ["雲量が多い", "最大仰角が低い"],
            message: "観測には厳しい条件です。",
        },
    ];

    // 3件の通過オブジェクトを作る モックデータ加工用関数
    const passes = mockConditions.map((condition) => {
        const startTime = now + condition.hoursLater * 60 * 60 * 1000;

        const maxTime = startTime + 3 * 60 * 1000;
        const endTime = startTime + 6 * 60 * 1000;

        return {
            id: `${cityId}_${Math.floor(startTime / 1000)}`,

            start: {
                time: new Date(startTime).toISOString(),
                direction: "北西",
                azimuth: 315,
                elevation: 10,
            },

            max: {
                time: new Date(maxTime).toISOString(),
                direction: "北",
                azimuth: 0,
                elevation: condition.elevation,
            },

            end: {
                time: new Date(endTime).toISOString(),
                direction: "北東",
                azimuth: 45,
                elevation: 10,
            },

            durationSec: 360,
            magnitude: condition.magnitude,
            cloudCover: condition.cloudCover,
            weatherAt: new Date(maxTime).toISOString(),
            visibility: condition.visibility,
            reasons: condition.reasons,
            message: condition.message,
        };
    });

    return {
        nextCandidate: passes[0],
        passes: passes,

        summary: {
            total: 3,
            good: 1,
            fair: 1,
            poor: 1,
            unknown: 0,
        },

        page: 1,
        limit: 3,
        total: 3,
        totalPages: 1,
    };
}

// ==================================================
// 通過1件の詳細を取得
// ==================================================

async function fetchISSPassDetail(cityId, passId) {
    const params = new URLSearchParams({
        city: cityId,
    });

    const response = await fetch(
        `${API_BASE_URL}/api/iss/passes/${encodeURIComponent(passId)}?${params}`,
    );

    if (!response.ok) {
        if (response.status === 404) {
            throw new Error(
                "通過情報が見つかりません。一覧を更新してください。",
            );
        }

        throw new Error(
            `通過詳細の取得に失敗しました: HTTP ${response.status}`,
        );
    }

    const data = await response.json();

    return data;
}

// ==================================================
// 3. 次回通過の表示
// ==================================================
// 通過開始までの残り時間を表示する
function updatePassCountdown(startTime) {
    const remainingSeconds = Math.max(
        0,
        Math.ceil((startTime - Date.now()) / 1000),
    );

    if (remainingSeconds === 0) {
        // タイマーが動いていたかを保存
        const wasRunning = countdownTimer !== null;

        // 今のカウントダウンを止める
        clearInterval(countdownTimer);
        countdownTimer = null;

        nextPassCountdown.textContent = "通過開始時刻です";

        // カウントダウン中に0になったときだけ再取得
        if (wasRunning) {
            refreshPasses();
        }

        return;
    }

    const hours = Math.floor(remainingSeconds / 3600);
    const minutes = Math.floor((remainingSeconds % 3600) / 60);
    const seconds = remainingSeconds % 60;
    //  先頭0うめ
    const pad = (value) => String(value).padStart(2, "0");

    nextPassCountdown.textContent = `${pad(hours)}時間 ${pad(minutes)}分 ${pad(seconds)}秒`;
}

function renderNextPass(pass) {
    // 都市変更時などに、前のタイマーを止める
    clearInterval(countdownTimer);
    countdownTimer = null;

    if (pass == null) {
        nextPassCountdown.textContent = "通過予定なし";
        nextPassStart.textContent = "—";
        nextPassDirection.textContent = "—";
        nextPassElevation.textContent = "—";
        nextPassCloud.textContent = "—";
        nextPassVisibility.textContent = "—";
        nextPassMessage.textContent = "今後3日間の通過予測はありません。";
        return;
    }

    // 日時を日本時間の表示にする
    nextPassStart.textContent = new Date(pass.start.time).toLocaleString(
        "ja-JP",
        {
            timeZone: "Asia/Tokyo",
            month: "2-digit",
            day: "2-digit",
            hour: "2-digit",
            minute: "2-digit",
            hour12: false,
        },
    );

    nextPassDirection.textContent = `${pass.start.direction} -> ${pass.end.direction}`;
    nextPassElevation.textContent = `${pass.max.elevation}°`;

    // 雲量の予報がない場合は不明
    nextPassCloud.textContent =
        pass.cloudCover === null ? "予報なし" : `${pass.cloudCover}%`;

    // GOから返る判定結果を日本語にする
    const visibilityLabels = {
        good: "良好",
        fair: "通常",
        poor: "条件が厳しい",
        unknown: "判定不明",
    };

    nextPassVisibility.textContent =
        visibilityLabels[pass.visibility] ?? "判定不明";

    nextPassMessage.textContent = pass.message;

    const startTime = new Date(pass.start.time).getTime();

    if (!Number.isFinite(startTime)) {
        nextPassCountdown.textContent = "開始時刻が不明です";
        return;
    }

    // 最初の表示
    updatePassCountdown(startTime);

    // 未来の通過なら、1秒ごとに表示を更新する
    if (startTime > Date.now()) {
        countdownTimer = setInterval(() => {
            updatePassCountdown(startTime);
        }, 1000);
    }
}

// ==================================================
// 通過一覧の表示
// ==================================================
function renderPassList(passes) {
    // 前の一覧や案内文を消す
    passesList.replaceChildren();
    //なかったらPタグつくってメッセージ出す
    if (passes.length === 0) {
        const message = document.createElement("p");

        message.className = "passes-empty";
        message.textContent = "表示する通過予定がありません。";

        passesList.append(message);
        return;
    }

    const visibilityLabels = {
        good: "良好",
        fair: "通常",
        poor: "条件が厳しい",
        unknown: "判定不明",
    };

    for (const pass of passes) {
        // 通過1件を表示する行
        const row = document.createElement("div");
        row.className = "pass-row";

        // 通過開始日時
        const date = document.createElement("span");
        date.className = "pass-row-date";
        date.textContent = new Date(pass.start.time).toLocaleString("ja-JP", {
            timeZone: "Asia/Tokyo",
            month: "2-digit",
            day: "2-digit",
            hour: "2-digit",
            minute: "2-digit",
            hour12: false,
        });

        // 方角
        const direction = document.createElement("span");
        direction.className = "pass-row-direction";
        direction.textContent = `${pass.start.direction} → ${pass.end.direction}`;

        // 最大仰角
        const elevation = document.createElement("span");
        elevation.className = "pass-row-elevation";
        elevation.textContent = `${pass.max.elevation}°`;

        // 見やすさ
        const visibility = document.createElement("span");
        visibility.className = "visibility-badge";
        visibility.dataset.visibility = pass.visibility;
        visibility.textContent =
            visibilityLabels[pass.visibility] ?? "判定不明";

        // 詳細ボタン
        const detailButton = document.createElement("button");
        detailButton.type = "button";
        detailButton.className = "pass-detail-button";
        detailButton.textContent = "詳細";
        detailButton.setAttribute("aria-expanded", "false");

        // 詳細を表示する領域
        const detailPanel = document.createElement("div");
        detailPanel.className = "pass-detail";
        detailPanel.id = `pass-detail-${pass.id}`;
        detailPanel.hidden = true;

        detailButton.setAttribute("aria-controls", detailPanel.id);

        // クリック時に使う都市を保存
        const cityId = citySelect.value;

        // 一度取得した詳細は、開き直しても再取得しない
        let detailLoaded = false;

        detailButton.addEventListener("click", async () => {
            // 開いている場合は閉じる
            if (!detailPanel.hidden) {
                detailPanel.hidden = true;
                detailButton.textContent = "詳細";
                detailButton.setAttribute("aria-expanded", "false");
                return;
            }

            // 閉じている場合は開く
            detailPanel.hidden = false;
            detailButton.textContent = "閉じる";
            detailButton.setAttribute("aria-expanded", "true");

            if (detailLoaded) return;

            detailPanel.textContent = "詳細を取得しています…";
            detailButton.disabled = true;

            try {
                const detail = await fetchISSPassDetail(cityId, pass.id);

                renderPassDetail(detailPanel, detail);
                detailLoaded = true;
            } catch (error) {
                console.error("通過詳細の取得に失敗しました:", error);

                detailPanel.textContent =
                    "通過詳細を取得できませんでした。通信状況やAPIの起動状態を確認してください。";
            } finally {
                detailButton.disabled = false;
            }
        });

        // 各項目を行へ入れる
        row.append(date, direction, elevation, visibility, detailButton);

        // 行と詳細をまとめる
        const item = document.createElement("div");
        item.className = "pass-item";
        item.append(row, detailPanel);

        // 一覧へ入れる
        passesList.append(item);
    }
}
// ==================================================
// 通過詳細の表示
// ==================================================

function renderPassDetail(container, pass) {
    container.replaceChildren();

    function formatTime(value) {
        return new Date(value).toLocaleTimeString("ja-JP", {
            timeZone: "Asia/Tokyo",
            hour: "2-digit",
            minute: "2-digit",
            second: "2-digit",
            hour12: false,
        });
    }

    const cloud = pass.cloudCover === null ? "予報なし" : `${pass.cloudCover}%`;

    const magnitude =
        pass.magnitude === null ? "不明" : `${pass.magnitude}等級`;

    const lines = [
        `開始：${formatTime(pass.start.time)}（${pass.start.direction}）`,
        `最高到達：${formatTime(pass.max.time)}（${pass.max.direction}）`,
        `終了：${formatTime(pass.end.time)}（${pass.end.direction}）`,
        `時刻は日本時間（JST）です。`,
        `最大仰角：${pass.max.elevation}°`,
        `目視可能時間：${pass.durationSec}秒`,
        `雲量：${cloud}`,
        `明るさ：${magnitude}`,
        `判定理由：${pass.reasons.join("・")}`,
        pass.message,
    ];

    for (const text of lines) {
        const paragraph = document.createElement("p");
        paragraph.textContent = text;
        container.append(paragraph);
    }
}

// 見やすさ別の件数
function renderPassSummary(summary) {
    passesSummary.textContent =
        `良好 ${summary.good}件 ／ ` +
        `通常 ${summary.fair}件 ／ ` +
        `条件が厳しい ${summary.poor}件 ／ ` +
        `判定不明 ${summary.unknown}件`;
}

// ページ情報とボタンの状態
function renderPassPagination(data) {
    if (data.total === 0) {
        passesPageInfo.textContent = "0件";
    } else if (data.passes.length === 0) {
        passesPageInfo.textContent = `全${data.total}件 ／ このページのデータはありません`;
    } else {
        const first = (data.page - 1) * data.limit + 1;
        const last = first + data.passes.length - 1;

        passesPageInfo.textContent =
            `全${data.total}件中 ${first}〜${last}件` +
            ` ／ ${data.page} / ${data.totalPages}ページ`;
    }

    passesPrev.disabled = data.page <= 1;
    passesNext.disabled = data.page >= data.totalPages;
}

// ==================================================
// 4. 通過情報の取得と表示
// ==================================================

async function refreshPasses() {
    const requestId = ++passesRequestId;
    const cityId = citySelect.value;
    const cityName = citySelect.selectedOptions[0].textContent;

    // HTMLで選んだ条件をまとめる
    const options = {
        visibility: visibilityFilter.value,
        minElevation: Number(elevationFilter.value),
        sort: passesSort.value,
        page: currentPage,
        limit: 3,
    };

    // 前の都市の情報とタイマーをクリア
    renderNextPass(null);
    nextPassCountdown.textContent = "取得中…";
    nextPassCity.textContent = `観測地点：${cityName}`;
    nextPassMessage.textContent = "通過情報を取得しています…";

    passesList.textContent = "通過情報を取得しています…";
    passesSummary.textContent = "";
    passesPageInfo.textContent = "";

    passesPrev.disabled = true;
    passesNext.disabled = true;

    try {
        const data = await fetchISSPasses(cityId, options);

        // より新しいリクエストがあれば、この結果は使わない
        if (requestId !== passesRequestId) return;

        currentPage = data.page;

        renderNextPass(data.nextCandidate);
        renderPassList(data.passes);
        renderPassSummary(data.summary);
        renderPassPagination(data);
    } catch (error) {
        if (requestId !== passesRequestId) return;

        console.error("通過情報の取得に失敗しました:", error);

        nextPassCountdown.textContent = "取得失敗";
        nextPassMessage.textContent = "通過情報を取得できませんでした。";

        passesList.textContent = "通過情報を取得できませんでした。";
    }
}

// ==================================================
// 5. 初期表示とイベント登録
// ==================================================
// 都市・絞り込み・並べ替えを変えたら1ページ目に戻す
function resetPageAndRefresh() {
    currentPage = 1;
    refreshPasses();
}
citySelect.addEventListener("change", resetPageAndRefresh);
visibilityFilter.addEventListener("change", resetPageAndRefresh);
elevationFilter.addEventListener("change", resetPageAndRefresh);
passesSort.addEventListener("change", resetPageAndRefresh);

// 前のページ
passesPrev.addEventListener("click", () => {
    if (currentPage > 1) {
        currentPage -= 1;
        refreshPasses();
    }
});

// 次のページ
passesNext.addEventListener("click", () => {
    currentPage += 1;
    refreshPasses();
});

// 都市の準備が終わってから通過情報を取得する
locationsReady.then((ready) => {
    if (ready) {
        refreshPasses();
        return;
    }

    renderNextPass(null);

    nextPassCountdown.textContent = "取得できません";
    nextPassMessage.textContent =
        "都市一覧を取得できないため、通過情報を表示できません。";

    passesList.textContent =
        "都市一覧を取得できないため、通過情報を表示できません。";
});
