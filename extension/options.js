function saveOptions(e) {
    e.preventDefault();
    browser.storage.sync.set({
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

    let getting = browser.storage.sync.get("backend");
    getting.then(setCurrentChoice, onError);
}

document.addEventListener("DOMContentLoaded", restoreOptions);
document.querySelector("form").addEventListener("submit", saveOptions);
