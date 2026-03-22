async function createShort() {
    const url = document.getElementById("url-input").value;
    const custom = document.getElementById("custom-input").value;

    clearCreate();

    const res = await fetch("/shorten", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            url: url,
            short_url: custom || undefined
        })
    });

    const data = await res.json();

    if (!res.ok) {
        document.getElementById("create-error").innerText = data.error || "Ошибка";
        return;
    }

    const shortUrl = `${window.location.origin}/s/${data.result.short_url}`;

    document.getElementById("short-result").innerHTML =
        `Короткая ссылка: <a href="${shortUrl}" target="_blank">${shortUrl}</a>`;
}

async function loadAnalytics() {
    const short = document.getElementById("short-input").value;
    const aggregate = document.getElementById("aggregate-select").value;

    clearAnalytics();

    let url = `/analytics/${short}`;
    if (aggregate) {
        url += `?aggregate_by=${aggregate}`;
    }

    const res = await fetch(url);
    const data = await res.json();

    if (!res.ok) {
        document.getElementById("analytics-error").innerText = data.error || "Ошибка";
        return;
    }

    renderTable(data);
}

function renderTable(data) {
    const tableHead = document.querySelector("#analytics-table thead");
    const tableBody = document.querySelector("#analytics-table tbody");

    tableHead.innerHTML = "";
    tableBody.innerHTML = "";

    const statsRaw = data.result?.stats || data.stats || [];

    let stats;
    if (Array.isArray(statsRaw)) {
        stats = statsRaw;
    } else {
        stats = [statsRaw];
    }

    if (stats.length === 0) {
        tableBody.innerHTML = "<tr><td>Нет данных</td></tr>";
        return;
    }

    const headers = Object.keys(stats[0]);

    tableHead.innerHTML =
        "<tr>" + headers.map(h => `<th>${h}</th>`).join("") + "</tr>";

    stats.forEach(row => {
        const tr = document.createElement("tr");
        tr.innerHTML = headers
            .map(h => `<td>${formatValue(row[h])}</td>`)
            .join("");
        tableBody.appendChild(tr);
    });
}

function formatValue(value) {
    if (Array.isArray(value) || typeof value === "object") {
        return `<pre>${JSON.stringify(value, null, 2)}</pre>`;
    }
    return value;
}

function clearCreate() {
    document.getElementById("create-error").innerText = "";
    document.getElementById("short-result").innerText = "";
}

function clearAnalytics() {
    document.getElementById("analytics-error").innerText = "";
    document.querySelector("#analytics-table thead").innerHTML = "";
    document.querySelector("#analytics-table tbody").innerHTML = "";
}