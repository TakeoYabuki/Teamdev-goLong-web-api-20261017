// ==================================================
// 通過情報の表示
// app.js の後に読み込むこと。
// citySelect・locationsReady・API_BASE_URL を使用する。
// ==================================================

// ==================================================
// 1. HTML要素の取得
// ==================================================
//左側：次回通過
const nextPassCity = document.getElementById("next-pass-city");
const nextPassCountdown = document.getElementById("next-pass-countdown");
const nextPassStart = document.getElementById("next-pass-start");
const nextPassDirection = document.getElementById("next-pass-direction");
const nextPassElevation = document.getElementById("next-pass-elevation");
const nextPassCloud = document.getElementById("next-pass-cloud");
const nextPassVisibility = document.getElementById("next-pass-visibility");
const nextPassMessage = document.getElementById("next-pass-message");
// 右側：通過一覧
const passesList = document.getElementById("passes-list");
const passesSummary = document.getElementById("passes-summary");
const passesPageInfo = document.getElementById("passes-page-info");
const passesPrev = document.getElementById("passes-prev");
const passesNext = document.getElementById("passes-next");
// 絞り込みと並べ替え
const visibilityFilter = document.getElementById("visibility-filter");
const elevationFilter = document.getElementById("elevation-filter");
const passesSort = document.getElementById("passes-sort");

// 1ページに表示する件数
const PASSES_PER_PAGE = 3;

// API完成後は false に変更する。
// 一覧と詳細を同時に切り替える。
const USE_MOCK_PASSES = false;

// 現在のページ
let currentPage = 1;

// 古いリクエストの結果で画面が上書きされないようにする
let passesRequestId = 0;

// カウントダウンのタイマーを管理
let countdownTimer = null;

// ==================================================
// 2. 共通の定数と表示用の関数
// ==================================================

// APIから返る見やすさの判定を日本語にする
const VISIBILITY_LABELS = {
    good: "良好",
    fair: "通常",
    poor: "条件が厳しい",
    unknown: "判定不明",
};

function formatVisibility(visibility) {
    return VISIBILITY_LABELS[visibility] ?? "判定不明";
}

// 「月/日 時：分 （日本時間）
function formatDateTime(value) {
    return new Date(value).toLocaleString("ja-JP",
        {
            timeZone: "Asia/Tokyo",
            month: "2-digit",
            day: "2-digit",
            hour: "2-digit",
            minute: "2-digit",
            hour12: false,
        }
    );
}

// 「時:分:秒」（日本時間）
function formatTime(value) {
    return new Date(value).toLocaleTimeString("ja-JP", {
        timeZone: "Asia/Tokyo",
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit",
        hour12: false,
    });
}

// 「北西 → 北東」
function formatDirection(pass) {
    return `${pass.start.direction} → ${pass.end.direction}`;
}

// 雲量の予報がない場合は「予報なし」
function formatCloud(cloudCover) {
    return cloudCover === null ? "予報なし" : `${cloudCover}%`;
}

// 取得中・失敗・都市なしなどのとき、左右の表示をまとめて切り替える
function showPassesMessage(countdownText, message) {
    clearNextPassDisplay();

    nextPassCountdown.textContent = countdownText;
    nextPassMessage.textContent = message;

    passesList.textContent = message;
    passesSummary.textContent = "";
    passesPageInfo.textContent = "";

    passesPrev.disabled = true;
    passesNext.disabled = true;
}
// ==================================================
// 3. 通過情報の取得（現在はモック）
// ==================================================
// モックの基準時刻を固定する。
// 再取得しても、同じ通過のIDが変わらないようにする。
const MOCK_BASE_TIME = Math.floor(Date.now() / 1000) * 1000;

// GET /api/iss/passes
async function fetchISSPasses(cityId, options = {}) {
    const {
        visibility = "",
        minElevation = 0,
        sort = "start",
        page = 1,
        limit = 3,
    } = options;

    if (USE_MOCK_PASSES) {
        return createMockPassesResponse(cityId, options);
    }

    // 都市IDをクエリパラメータにする
    const params = new URLSearchParams({
        city: cityId,
        minElevation: String(minElevation),
        sort: sort,
        page: String(page),
        limit: String(limit),
    });

    if (visibility !== "") {
        params.set("visibility", visibility);
    }
    return fetchPassesJSON(
        `${API_BASE_URL}/api/iss/passes?${params}`,
        "通過情報を取得できませんでした。",
    );
}

// HTTPエラーと、JSONとして読めない応答を確認する
async function fetchPassesJSON(url, fallbackMessage) {
    const response = await fetch(url);
    let data;

    try {
        data = await response.json();
    } catch {
        const error = new Error(
            response.ok
                ? "APIの応答をJSONとして読み取れませんでした。"
                : fallbackMessage,
        );

        error.status = response.status;
        throw error;
    }

    if (!response.ok) {
        // Go側の共通エラー形式から取り出す
        const message = data?.error?.message;

        const error = new Error(
            typeof message === "string" && message.trim() !== ""
                ? message
                : fallbackMessage,
        );

        error.status = response.status;

        if (typeof data?.error?.code === "string") {
            error.code = data.error.code;
        }

        throw error;
    }

    return data;
}

// ==================================================
// モックデータ
// ==================================================

// 一覧と詳細で同じID・同じ内容の通過を作る
function createMockPasses(cityId) {
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

    return mockConditions.map((condition) => {
        const startTime =
            MOCK_BASE_TIME + condition.hoursLater * 60 * 60 * 1000;

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
}
function createMockPassesResponse(cityId, options) {
    const { page = 1, limit = PASSES_PER_PAGE } = options;

    const from = Date.now();
    const to = from + 72 * 60 * 60 * 1000;

    // 開始時刻を過ぎた通過は次回候補と一覧から外す
    const passes = createMockPasses(cityId).filter((pass) => {
        const start = Date.parse(pass.start.time);

        return start >= from && start < to;
    });

    const summary = {
        total: passes.length,
        good: 0,
        fair: 0,
        poor: 0,
        unknown: 0,
    };

    for (const pass of passes) {
        summary[pass.visibility] += 1;
    }

    // モックは絞り込み・並べ替えを再現しない。
    // ページ分けだけ確認できる。
    return {
        nextCandidate: passes[0] ?? null,
        passes: passes.slice((page - 1) * limit, page * limit),
        summary,
        page,
        limit,
        total: passes.length,
        totalPages: Math.ceil(passes.length / limit),
    };
}
// ==================================================
// 4. 通過1件の詳細を取得
// ==================================================

// GET /api/iss/passes/{id}
async function fetchISSPassDetail(cityId, passId) {
    if (USE_MOCK_PASSES) {
        const pass = createMockPasses(cityId).find(
            (item) => item.id === passId,
        );

        if (!pass) {
            const error = new Error(
                "詳細を取得できません。一覧を更新してください。",
            );

            error.status = 404;
            throw error;
        }

        return pass;
    }
    const params = new URLSearchParams({
        city: cityId,
    });

    return fetchPassesJSON(
        `${API_BASE_URL}/api/iss/passes/${encodeURIComponent(passId)}?${params}`,
        "通過詳細を取得できませんでした。",
    );
}

// ==================================================
// 5. 次回通過の表示
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

// タイマーの停止と数値欄の初期化だけを担当する関数
function clearNextPassDisplay() {
    clearInterval(countdownTimer);
    countdownTimer = null;

    nextPassCountdown.textContent = "—";
    nextPassStart.textContent = "—";
    nextPassDirection.textContent = "—";
    nextPassElevation.textContent = "—";
    nextPassCloud.textContent = "—";
    nextPassVisibility.textContent = "—";
    nextPassMessage.textContent = "";
}

function renderNextPass(pass) {
    clearNextPassDisplay();

    if (pass === null) {
        nextPassCountdown.textContent = "通過予定なし";
        nextPassMessage.textContent = "今後3日間の通過予測はありません。";
        return;
    }

    nextPassStart.textContent = formatDateTime(pass.start.time);
    nextPassDirection.textContent = formatDirection(pass);
    nextPassElevation.textContent = `${pass.max.elevation}°`;
    nextPassCloud.textContent = formatCloud(pass.cloudCover);
    nextPassVisibility.textContent = formatVisibility(pass.visibility);
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
// 6. 通過一覧の表示
// ==================================================
function renderPassList(passes, summary, total) {
    // 前の一覧や案内文を消す
    passesList.replaceChildren();
    //なかったらPタグつくってメッセージ出す
    if (passes.length === 0) {
        const message = document.createElement("p");
        message.className = "passes-empty";

        if (summary.total === 0) {
            message.textContent = "今後3日間に通過予定はありません。";
        } else if (total === 0) {
            message.textContent = "指定した条件に合う通過予定はありません。";
        } else {
            message.textContent = "このページに通過情報はありません。";
        }

        passesList.append(message);
        return;
    }

    // 描画した時点の都市を詳細取得に使う
    const cityId = citySelect.value;

    for (const pass of passes) {
        passesList.append(createPassItem(pass, cityId));
    }
}

// 通過1件分の行と詳細パネルを作る
function createPassItem(pass, cityId) {
    const row = document.createElement("div");
    row.className = "pass-row";

    // 通過開始日時
    const date = document.createElement("span");
    date.className = "pass-row-date";
    date.textContent = formatDateTime(pass.start.time);

    // 方角
    const direction = document.createElement("span");
    direction.className = "pass-row-direction";
    direction.textContent = formatDirection(pass);

    // 最大仰角
    const elevation = document.createElement("span");
    elevation.className = "pass-row-elevation";
    elevation.textContent = `${pass.max.elevation}°`;

    // 見やすさ
    const visibility = document.createElement("span");
    visibility.className = "visibility-badge";
    visibility.dataset.visibility = pass.visibility;
    visibility.textContent = formatVisibility(pass.visibility);

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

            validatePass(detail, "通過詳細");
            renderPassDetail(detailPanel, detail);

            detailLoaded = true;
        } catch (error) {
            console.error("通過詳細の取得に失敗しました:", error);

            detailPanel.textContent =
                error.status === 404
                    ? "詳細を取得できません。一覧を更新してください。"
                    : "通過詳細を取得できませんでした。通信状況やAPIの起動状態を確認してください。";
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

    return item;
}

// ==================================================
// 7. 通過詳細の表示
// ==================================================

function renderPassDetail(container, pass) {
    container.replaceChildren();

    const magnitude =
        pass.magnitude === null ? "不明" : `${pass.magnitude}等級`;

    const lines = [
        `開始：${formatTime(pass.start.time)}（${pass.start.direction}）`,
        `最高到達：${formatTime(pass.max.time)}（${pass.max.direction}）`,
        `終了：${formatTime(pass.end.time)}（${pass.end.direction}）`,
        `時刻は日本時間（JST）です。`,
        `最大仰角：${pass.max.elevation}°`,
        `目視可能時間：${pass.durationSec}秒`,
        `雲量：${formatCloud(pass.cloudCover)}`,
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
        `${VISIBILITY_LABELS.good} ${summary.good}件 ／ ` +
        `${VISIBILITY_LABELS.fair} ${summary.fair}件 ／ ` +
        `${VISIBILITY_LABELS.poor} ${summary.poor}件 ／ ` +
        `${VISIBILITY_LABELS.unknown} ${summary.unknown}件`;
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
// 8. 応答の検証
// ==================================================
// 通過1件の内容を確認する。
// nullを許可するかどうかは、呼び出し側で判断する。
function validatePass(pass, label = "通過情報") {
    const isObject = (value) =>
        value !== null && typeof value === "object" && !Array.isArray(value);

    const isText = (value) => typeof value === "string" && value.trim() !== "";

    if (!isObject(pass) || !isText(pass.id)) {
        throw new Error(`${label}のIDまたは形式が不正です。`);
    }

    for (const key of ["start", "max", "end"]) {
        const point = pass[key];

        if (
            !isObject(point) ||
            !isText(point.time) ||
            !Number.isFinite(Date.parse(point.time)) ||
            !isText(point.direction) ||
            !Number.isFinite(point.azimuth) ||
            point.azimuth < 0 ||
            point.azimuth >= 360 ||
            !Number.isFinite(point.elevation) ||
            point.elevation < 0 ||
            point.elevation > 90
        ) {
            throw new Error(`${label}の${key}が不正です。`);
        }
    }

    if (
        !Number.isInteger(pass.durationSec) ||
        pass.durationSec < 0 ||
        (pass.cloudCover !== null &&
            (!Number.isFinite(pass.cloudCover) ||
                pass.cloudCover < 0 ||
                pass.cloudCover > 100)) ||
        (pass.magnitude !== null && !Number.isFinite(pass.magnitude)) ||
        !Object.hasOwn(VISIBILITY_LABELS, pass.visibility) ||
        !Array.isArray(pass.reasons) ||
        !pass.reasons.every((reason) => typeof reason === "string") ||
        typeof pass.message !== "string" ||
        (pass.weatherAt !== null &&
            (typeof pass.weatherAt !== "string" ||
                !Number.isFinite(Date.parse(pass.weatherAt))))
    ) {
        throw new Error(`${label}の表示項目が不正です。`);
    }
}
// API応答に、必要な項目があるか確認
function validatePassesResponse(data) {
    if (data === null || typeof data !== "object" || Array.isArray(data)) {
        throw new Error("通過情報の応答がオブジェクトではありません");
    }
    // nullは「候補なし」なので許可する
    if (!Object.hasOwn(data, "nextCandidate")) {
        throw new Error("次回候補の項目がありません");
    }
    if (
        data.nextCandidate !== null &&
        (typeof data.nextCandidate !== "object" ||
            Array.isArray(data.nextCandidate))
    ) {
        throw new Error("次回候補の形式が不正です");
    }
    if (!Array.isArray(data.passes)) {
        throw new Error("通過一覧が配列ではありません");
    }
    if (data.nextCandidate !== null) {
        validatePass(data.nextCandidate, "次回候補");
    }

    for (const pass of data.passes) {
        validatePass(pass, "通過一覧");
    }

    // 集計に必要な件数を確認
    const summaryKeys = ["total", "good", "fair", "poor", "unknown"];

    for (const key of summaryKeys) {
        const value = data.summary?.[key];

        if (!Number.isInteger(value) || value < 0) {
            throw new Error(`集計の${key}が不正です`);
        }
    }

    // ページ情報を確認
    const minimums = {
        page: 1,
        limit: 1,
        total: 0,
        totalPages: 0,
    };

    for (const [key, minimum] of Object.entries(minimums)) {
        if (!Number.isInteger(data[key]) || data[key] < minimum) {
            throw new Error(`ページ情報の${key}が不正です`);
        }
    }

    if (data.limit > 50) {
        throw new Error("1ページの取得件数が上限を超えています");
    }
}

// ==================================================
// 9. 通過情報の取得と表示
// ==================================================

async function refreshPasses() {
    const requestId = ++passesRequestId;
    const cityId = citySelect.value;
    const cityName = citySelect.selectedOptions[0]?.textContent ?? "";

    // HTMLで選んだ条件をまとめる
    const options = {
        visibility: visibilityFilter.value,
        minElevation: Number(elevationFilter.value),
        sort: passesSort.value,
        page: currentPage,
        limit: PASSES_PER_PAGE,
    };
    showPassesMessage("取得中...", "通過情報を取得しています...");

    nextPassCity.textContent = `観測地点：${cityName}`;

    try {
        const data = await fetchISSPasses(cityId, options);

        // より新しいリクエストがあれば、この結果は使わない
        if (requestId !== passesRequestId) return;

        // 追加：応答の形式を確認
        validatePassesResponse(data);

        currentPage = data.page;

        renderNextPass(data.nextCandidate);
        renderPassList(data.passes, data.summary, data.total);
        renderPassSummary(data.summary);
        renderPassPagination(data);
    } catch (error) {
        if (requestId !== passesRequestId) return;

        console.error("通過情報の取得に失敗しました", error);

        showPassesMessage("取得失敗", "通過情報を取得できませんでした。");
    }
}

// ==================================================
// 10. 初期表示とイベント登録
// ==================================================
// 都市・絞り込み・並べ替えを変えたら1ページ目に戻す
function resetPageAndRefresh() {
    currentPage = 1;
    refreshPasses();
}

// モックでは条件処理を実装していないため、選択欄を無効にする。
// 本番ではGo側が絞り込み・並べ替えを担当する。
if (USE_MOCK_PASSES) {
    for (const control of [visibilityFilter, elevationFilter, passesSort]) {
        control.disabled = true;
        control.title = "モック表示中は絞り込み・並べ替えを利用できません。";
    }
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

    // 前の都市の情報とタイマーをクリア
    clearNextPassDisplay();

    showPassesMessage(
        "取得できません",
        "都市一覧を取得できないため、通過情報を表示できません。",
    );
});
