// ==================================================
// 1. 設定
// ==================================================

const API_BASE_URL = "http://localhost:8000";

const NIGHT_TILE_URL =
    "https://gibs.earthdata.nasa.gov/wmts/epsg3857/best/" +
    "VIIRS_CityLights_2012/default/" +
    "GoogleMapsCompatible_Level8/{z}/{y}/{x}.jpeg";

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

// ==================================================
// 3. 地図の作成
// ==================================================

function createMap() {
    const map = L.map("map").setView([20, 0], 2); // Leaflet Map参照

    // 夜の世界地図を表示 TODO:出典明記すること
    L.tileLayer(NIGHT_TILE_URL, {
        attribution:
            'Imagery: <a href="https://earthdata.nasa.gov/gibs">NASA GIBS</a>' +
            " / Suomi NPP VIIRS",
        minZoom: 1,
        maxNativeZoom: 8,
        maxZoom: 8,
    }).addTo(map);

    return map;
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

// ==================================================
// 4. 観測都市の表示
// ==================================================
// TODO: 選択した都市の更新 cityMarkerとissMarkerに入れる zoomしたほうがいいかも
function updateSelectedCity() {
    const city = cities[citySelect.value]; // optionの選択肢を拾う

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

    cityMarker.setPopupContent(`観測地点: ${city.name}`); // cityMarkerでポップアップする
    selectedCity.textContent = `観測地点: ${city.name}`; // Dom追加
}

// ==================================================
// 5. ISSのデータ取得
// ==================================================

// 非同期 ISSの位置をapiから取ってくる
async function fetchIssPosition() {
    // TODO: nowが返す構造
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

    return data; // TODO: API完成したらこっちに変える
}

// ==================================================
// 6. ISSの表示
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
            fillColor: "#38d9ff",
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

    if (Number.isNaN(positionAt.getTime())) {
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

// ==================================================
// 7. 通信と表示をつなぐ処理
// ==================================================

async function refreshISS() {
    try {
        const position = await fetchIssPosition();

        updateISSMarker(position);
        updateISSStatus(position);
    } catch (error) {
        console.error("ISSの更新に失敗しました", error);
        // ISSを取得できなかった時のエラーハンドリング DOMに表示するほう
        issStatus.textContent =
            issMarker === null
                ? "ISSの位置を取得できませんでした。"
                : "ISSの更新に失敗しました。マーカーは前回取得した位置です。";
    }
}

// ==================================================
// 都市一覧の取得
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
        if (cities.tokyo && DEMO_CITY_IDS.includes("tokyo")) {
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
// 8. 初期表示とイベント登録
// ==================================================

// 3のmap展開
const map = createMap();

// 都市一覧の取得と初期表示を開始する
const locationsReady = initializeLocations();

// 初期表示
// updateSelectedCity();

// 都市の選択
citySelect.addEventListener("change", updateSelectedCity);

// 地図操作のボタン
focusCityButton.addEventListener("click", focusSelectedCity);
showWorldButton.addEventListener("click", showWorld);

// ISSの取得
async function startISSUpdates() {
    await refreshISS();
    setTimeout(startISSUpdates, 5000); // 5秒おきに繰り返す
}
issStatus.textContent = "ISSの位置を取得しています...";
startISSUpdates();
