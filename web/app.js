// ==================================================
// 1. 設定
// ==================================================

const API_BASE_URL = "http://localhost:8000";

const NIGHT_TILE_URL =
    "https://gibs.earthdata.nasa.gov/wmts/epsg3857/best/" +
    "VIIRS_CityLights_2012/default/" +
    "GoogleMapsCompatible_Level8/{z}/{y}/{x}.jpeg";

// TODO: 都市検索APIから取得するデータに置き換える
const cities = {
    tokyo: {
        name: "東京都心",
        lat: 35.6812,
        lon: 139.7671,
    },
    osaka: {
        name: "大阪市",
        lat: 34.6937,
        lon: 135.5023,
    },
    kagoshima: {
        name: "鹿児島市",
        lat: 31.5966,
        lon: 130.5571,
    },
};

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
    map.setView([city.lat, city.lon], 7);
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
        throw new Error(`ISSの一取得に失敗しました: HTTP ${response.status}`);
    }

    const data = await response.json();

    // TODO: データがおかしいときのエラーハンドリング。

    return data;
}

// ==================================================
// 6. ISSの表示
// ==================================================

// refreshISSから呼び出す
// TODO: ISSの位置をLeafteのcircleMarkerで表示
function updateISSMarker(position) {}

// TODO: ISSのデータをDomに挿入 textContent
function updateISSStatus(position) {}

// ==================================================
// 7. 通信と表示をつなぐ処理
// ==================================================

async function refreshISS() {
    try {
        const position = await fetchIssPosition();

        updateISSMarker(position);
        updateISSStatus(position);
    } catch (error) {
        // TODO: ISSを取得できなかった時のエラーハンドリング DOMに表示するほう
    }
}

// ==================================================
// 8. 初期表示とイベント登録
// ==================================================

// 3のmap展開
const map = createMap();

// 初期表示
updateSelectedCity();

// 都市の選択
citySelect.addEventListener("change", updateSelectedCity);

// 地図操作のボタン
focusCityButton.addEventListener("click", focusSelectedCity);
showWorldButton.addEventListener("click", showWorld);

// ISSの取得
issStatus.textContent = "ISSの位置を取得しています_";
refreshISS();
