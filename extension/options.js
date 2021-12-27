function saveOptions(e) {
    e.preventDefault();
    chrome.storage.local.set({
        backend: document.querySelector("#backend").value
    });
}

function restoreOptions() {
    function setCurrentChoice(result) {
        document.querySelector("#backend").value = result.backend || "http://51.75.67.146:8080/";
    }

    function onError(error) {
        console.log(`Error: ${error}`);
    }

    chrome.storage.local.get("backend", setCurrentChoice);
}

document.addEventListener("DOMContentLoaded", restoreOptions);
document.querySelector("form").addEventListener("submit", saveOptions);
