// ==================================================
// 1. 設定
// ==================================================

const API_BASE_URL = "http://localhost:8000";

const NIGHT_TILE_URL =
    "https://gibs.earthdata.nasa.gov/wmts/epsg3857/best/" +
    "VIIRS_CityLights_2012/default/" +
    "GoogleMapsCompatible_Level8/{z}/{y}/{x}.jpeg";

// 初期表示（世界全体）の中心とズーム
const WORLD_VIEW = {
    center: [20, 0],
    zoom: 2,
};

// 都市へ移動するときのズーム
const CITY_ZOOM = 5;

// ISSのマーカーと軌跡の色
const ISS_COLOR = "#38d9ff";

// ISSの位置を取り直す間隔
const ISS_UPDATE_INTERVAL_MS = 5000;

// 移動履歴は直近40分を表示する
const ISS_HISTORY_DURATION_MS = 40 * 60 * 1000;

// APIから取得した都市データを入れる
let cities = {};

// デモアプリで選択できる12都市
const DEMO_CITY_IDS = [
    "sapporo",
    "aomori",
    "sendai",
    "tokyo",
    "niigata",
    "nagoya",
    "osaka",
    "hiroshima",
    "kochi",
    "fukuoka",
    "kagoshima",
    "naha",
];

// ==================================================
// 2. HTML要素とマーカーの管理
// ==================================================

const citySelect = document.getElementById("city-select");
const selectedCity = document.getElementById("selected-city");
const issStatus = document.getElementById("iss-status");

// 地図操作のボタン
const focusCityButton = document.getElementById("focus-city");
const showWorldButton = document.getElementById("show-world");

let cityMarker = null; //
let issMarker = null;

let issHistory = [];
let issTrail = null; // leafletのpolylineオブジェクト

// ==================================================
// 3. 地図の作成
// ==================================================

function createMap() {
    const newMap = L.map("map").setView(WORLD_VIEW.center, WORLD_VIEW.zoom); // Leaflet Map参照

    L.tileLayer(NIGHT_TILE_URL, {
        attribution:
            'Imagery: <a href="https://earthdata.nasa.gov/gibs">NASA GIBS</a>' +
            " / Suomi NPP VIIRS",
        minZoom: 1,
        maxNativeZoom: 8,
        maxZoom: 8,
    }).addTo(newMap);

    return newMap;
}

// 仰角10度以上の通過があり得る緯度範囲の目安
// ISS高度約420kmを仮定した概算
function addObservationGuides(map) {
    const guides = [64, -64].map((latitude) => {
        const label = latitude > 0 ? "北緯約64度" : "南緯約64度";

        const line = L.polyline([], {
            color: "#6e8c90",
            weight: 1.5,
            opacity: 0.5,
            dashArray: "1, 5",
            lineCap: "round",
        })
            .addTo(map)
            .bindTooltip(`${label}：仰角10度以上の通過範囲の目安`);

        return { latitude, line };
    });

    function updateGuideLines() {
        const bounds = map.getBounds();

        // 地図の左右端より少し外まで線を伸ばす
        const west = bounds.getWest() - 10;
        const east = bounds.getEast() + 10;

        for (const { latitude, line } of guides) {
            line.setLatLngs([
                [latitude, west],
                [latitude, east],
            ]);
        }
    }

    // 初期表示
    updateGuideLines();

    // 移動・ズーム・表示サイズの変更に合わせて更新
    map.on("moveend resize", updateGuideLines);
}

// 選択中の都市へ移動
function focusSelectedCity() {
    const city = cities[citySelect.value];

    if (!city) return;
    map.setView([city.lat, city.lon], 5);
}

// 初期の世界表示に戻す
function showWorld() {
    map.setView([20, 0], 2);
}

// 地図の表示状態に合わせて、操作ボタンの色を切り替える
function updateMapButtonState() {
    const city = cities[citySelect.value];

    // 初期表示より拡大していて、観測地点が画面内にあるか
    const isCityVisible = Boolean(
        city &&
        map.getZoom() > WORLD_VIEW.zoom &&
        map.getBounds().contains([city.lat, city.lon]),
    );

    focusCityButton.classList.toggle("button-primary", !isCityVisible);
    focusCityButton.classList.toggle("button-secondary", isCityVisible);

    showWorldButton.classList.toggle("button-primary", isCityVisible);
    showWorldButton.classList.toggle("button-secondary", !isCityVisible);
}
// ==================================================
// 4. 観測都市の表示
// ==================================================

function updateSelectedCity() {
    const city = cities[citySelect.value];

    if (!city) {
        selectedCity.textContent = "都市を選択してください";
        return;
    }

    const coordinates = [city.lat, city.lon]; // 座標

    if (cityMarker == null) {
        cityMarker = L.marker(coordinates).addTo(map).bindPopup("");
    } else {
        cityMarker.setLatLng(coordinates);
    }

    cityMarker.setPopupContent(`観測地点: ${city.name}`);
    selectedCity.textContent = `観測地点: ${city.name}`;

    updateMapButtonState();
}

// ==================================================
// 5. 都市一覧の取得
// ==================================================

async function fetchLocations() {
    const params = new URLSearchParams({
        limit: "50",
    });

    const response = await fetch(
        `${API_BASE_URL}/api/locations/search?${params}`,
    );

    const data = await response.json().catch(() => ({}));
    if (!response.ok) {
        throw new Error(
            data.error?.message ?? "都市一覧を取得できませんでした。",
        );
    }

    if (!Array.isArray(data.results)) {
        throw new Error("都市一覧の形式が不正です");
    }
    return data.results;
}

async function initializeLocations() {
    citySelect.disabled = true;
    selectedCity.textContent = "都市一覧を取得しています…";

    try {
        const locations = await fetchLocations();

        // APIの都市データを、IDで取り出せる形にする
        for (const location of locations) {
            cities[location.id] = location;
        }

        // HTMLに書いていた選択肢を削除する
        citySelect.replaceChildren();

        // デモで使う都市だけを選択肢に追加する
        for (const cityId of DEMO_CITY_IDS) {
            const city = cities[cityId];

            if (!city) continue;

            const option = document.createElement("option");
            option.value = city.id;
            option.textContent = city.name;

            citySelect.append(option);
        }

        if (citySelect.options.length === 0) {
            throw new Error("選択できる都市がありません");
        }

        // 初期選択は東京にする
        if (cities.tokyo) {
            citySelect.value = "tokyo";
        }

        updateSelectedCity();
        citySelect.disabled = false;

        return true;
    } catch (error) {
        console.error("都市一覧の取得に失敗しました", error);

        citySelect.replaceChildren();
        selectedCity.textContent =
            "都市一覧を取得できませんでした。APIの起動状態を確認してください。";

        return false;
    }
}

// ==================================================
// 6. ISSのデータ取得
// ==================================================

// 非同期 ISSの位置をapiから取ってくる
async function fetchIssPosition() {
    const response = await fetch(`${API_BASE_URL}/api/iss/now`);
    if (!response.ok) {
        throw new Error(`ISSの位置取得に失敗しました: HTTP ${response.status}`);
    }

    const data = await response.json();

    // データがおかしいときのエラーハンドリング。
    if (
        !Number.isFinite(data.lat) ||
        !Number.isFinite(data.lon) ||
        data.lat < -90 ||
        data.lat > 90 ||
        data.lon < -180 ||
        data.lon > 180
    ) {
        throw new Error("ISSの位置データが不正です");
    }

    return data;
}
// ==================================================
// 7. ISSの表示
// ==================================================

// refreshISSから呼び出す
// ISSの位置をLeafteのcircleMarkerで表示
function updateISSMarker(position) {
    const coordinates = [position.lat, position.lon];

    if (issMarker == null) {
        issMarker = L.circleMarker(coordinates, {
            radius: 8,
            color: "#ffffff",
            weight: 2,
            fillColor: ISS_COLOR,
            fillOpacity: 1,
        })
            .addTo(map)
            .bindPopup("ISSの現在位置");
    } else {
        issMarker.setLatLng(coordinates);
    }
}

// ISSのデータをDomに挿入 textContent
function updateISSStatus(position) {
    const positionAt = new Date(position.positionAt);

    if (
        typeof position.positionAt !== "string" ||
        Number.isNaN(positionAt.getTime())
    ) {
        issStatus.textContent =
            "ISSの位置を取得しました。データ時刻は不明です。";
        return;
    }

    const formattedTime = positionAt.toLocaleString("ja-JP", {
        timeZone: "Asia/Tokyo",
        hour12: false,
    });

    issStatus.textContent = `ISSのデータ時刻: ${formattedTime}(日本時間)`;
}

// 取得した位置を保存し、移動履歴の線を更新する
function updateISSTrail(position) {
    const positionTime = new Date(position.positionAt).getTime();

    // データ時刻が不正なら、履歴には追加しない
    if (
        typeof position.positionAt !== "string" ||
        !Number.isFinite(positionTime)
    ) {
        return;
    }

    // 同じ時刻や、前より古いデータは追加しない
    const lastPosition = issHistory[issHistory.length - 1];

    if (lastPosition && positionTime <= lastPosition.time) {
        return;
    }

    issHistory.push({
        lat: position.lat,
        lon: position.lon,
        time: positionTime,
    });

    // 設定した保持期間より古い位置を削除する
    const cutoff = positionTime - ISS_HISTORY_DURATION_MS;

    issHistory = issHistory.filter((point) => point.time >= cutoff);

    // 地図の左右の端をまたぐ場合は、線を分ける
    const segments = [];
    let segment = [];

    for (const point of issHistory) {
        const previous = segment[segment.length - 1];

        if (previous && Math.abs(point.lon - previous[1]) > 180) {
            segments.push(segment);
            segment = [];
        }

        segment.push([point.lat, point.lon]);
    }

    if (segment.length > 0) {
        segments.push(segment);
    }

    // 線は最初の1回だけ作成する
    if (issTrail === null) {
        issTrail = L.polyline([], {
            color: ISS_COLOR,
            weight: 3,
            opacity: 0.3,
            interactive: false,
        }).addTo(map);

        // 現在位置のマーカーより後ろに表示する
        issTrail.bringToBack();
    }
    // 保存した位置で線を更新する
    issTrail.setLatLngs(segments);
}

// ==================================================
// 8. 通信と表示をつなぐ処理
// ==================================================

async function refreshISS() {
    try {
        const position = await fetchIssPosition();

        updateISSMarker(position);
        updateISSStatus(position);
        updateISSTrail(position);
    } catch (error) {
        console.error("ISSの更新に失敗しました", error);
        // ISSを取得できなかった時のエラーハンドリング DOMに表示するほう
        issStatus.textContent =
            issMarker === null
                ? "ISSの位置を取得できませんでした。"
                : "ISSの更新に失敗しました。マーカーは前回取得した位置です。";
    }
}

// 取得が終わってから次の取得を予約する（リクエストが重ならない）
async function startISSUpdates() {
    await refreshISS();
    setTimeout(startISSUpdates, ISS_UPDATE_INTERVAL_MS);
}

// ==================================================
// 9. 初期表示とイベント登録
// ==================================================

// 3のmap展開
const map = createMap();
addObservationGuides(map);
// 手動操作・ボタン操作・画面サイズ変更に対応する
map.on("moveend zoomend resize", updateMapButtonState);

// 初期表示の色を設定する
updateMapButtonState();

// 都市一覧の取得と初期表示を開始する
const locationsReady = initializeLocations();

// 初期表示
// updateSelectedCity();

// 都市の選択
citySelect.addEventListener("change", updateSelectedCity);

// 地図操作のボタン
focusCityButton.addEventListener("click", focusSelectedCity);
showWorldButton.addEventListener("click", showWorld);

issStatus.textContent = "ISSの位置を取得しています...";
startISSUpdates();
