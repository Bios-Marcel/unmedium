document.addEventListener("DOMContentLoaded", function (event) {
    if (document.body.innerHTML.indexOf("<script src=\"https://cdn-client.medium.com/") != -1) {
        chrome.storage.local.get("backend", function (data) {
            window.location = (data.backend || "http://51.75.67.146:8080/") + window.location.href;
        });
    }
});