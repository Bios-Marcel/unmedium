document.addEventListener("DOMContentLoaded", function (event) {
    if (document.body.innerHTML.indexOf("<script src=\"https://cdn-client.medium.com/") != -1) {
        browser.storage.sync.get("backend").then(function (data) {
            window.location = data.backend + window.location.href;
        }, null);
    }
});