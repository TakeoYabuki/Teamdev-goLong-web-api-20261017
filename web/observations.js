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

// カウントダウンのタイマーを管理
let countdownTimer = null;

// ==================================================
// 2. 通過情報の取得（現在はモック）
// ==================================================
// GET /api/iss/passesを呼ぶ cityIDにはtokyoとかが入る
async function fetchISSPasses(cityId) {
    // 都市IDをクエリパラメータにする
    // const params = new URLSearchParams({
    //     city: cityId,
    // });

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

    const startTime = Date.now() + 2 * 60 * 60 * 1000;

    return {
        // TODO: APIができたらfetchに変える
        nextCandidate: {
            id: `${cityId}_${Math.floor(startTime / 1000)}`,

            start: {
                time: new Date(startTime).toISOString(),
                direction: "北西",
                azimuth: 315,
                elevation: 10,
            },

            max: {
                time: new Date(startTime + 3 * 60 * 1000).toISOString(),
                direction: "北",
                azimuth: 0,
                elevation: 62,
            },

            end: {
                time: new Date(startTime + 6 * 60 * 1000).toISOString(),
                direction: "北東",
                azimuth: 45,
                elevation: 10,
            },

            durationSec: 360,
            magnitude: -3,
            cloudCover: 20,
            weatherAt: new Date(startTime + 3 * 60 * 1000).toISOString(),
            visibility: "poor",
            reasons: ["雲量が少ない", "最大仰角が高い"],
            message: "空の高い位置を通り、観測しやすい条件です。",
        },
    };
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

    const hours = Math.floor(remainingSeconds / 3600);
    const minutes = Math.floor((remainingSeconds % 3600) / 60);
    const seconds = remainingSeconds % 60;

    // 先頭0うめ
    const pad = (value) => String(value).padStart(2, "0");
    nextPassCountdown.textContent = `${pad(hours)}時間 ${pad(minutes)}分 ${pad(seconds)}秒`;

    // 開始時間になったらタイマーを止める
    if (remainingSeconds === 0) {
        clearInterval(countdownTimer);
        countdownTimer = null;
        nextPassCountdown.textContent = "通過開始時刻です";
    }
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
        pass.cloudCover === null ? "不明" : `${pass.cloudCover}%`;

    // GOから返る判定結果を日本語にする
    const visibilityLabels = {
        good: "良好",
        fair: "普通",
        poor: "条件が厳しい",
        unknown: "判定不能",
    };

    nextPassVisibility.textContent =
        visibilityLabels[pass.visibility] ?? "判定不能";

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
// 4. 通過情報の取得と表示
// ==================================================

async function refreshPasses() {
    const cityId = citySelect.value; // htmlからIDを取得
    // 選択肢に表示されている都市名
    const cityName = citySelect.selectedOptions[0].textContent;

    nextPassCity.textContent = `観測地点: ${cityName}`;
    nextPassMessage.textContent = `通過情報を取得しています...`;

    try {
        const data = await fetchISSPasses(cityId);
        // 取得中に都市が変わった場合は表示しない
        if (citySelect.value !== cityId) return;

        // 取得した次回通過をカードに表示
        renderNextPass(data.nextCandidate);
    } catch (error) {
        if (citySelect.value !== cityId) return;

        console.error("通過情報の取得に失敗しました", error);

        renderNextPass(null);
        nextPassCountdown.textContent = "取得失敗";
        nextPassMessage.textContent = "通過情報を取得できませんでした。";
    }
}

// ==================================================
// 5. 初期表示とイベント登録
// ==================================================

// 都市変更時に、情報カードも更新する
citySelect.addEventListener("change", refreshPasses);

// 最初の1回
refreshPasses();
