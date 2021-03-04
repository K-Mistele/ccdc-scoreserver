window.addEventListener('load', () => {

    // PREVENT WHITESPACE FROM BEING INSERTED INTO THE SERVICE NAME FIELD
    const nameInput = document.getElementById('serviceName');
    nameInput.addEventListener('keypress', (event) => {
        const key = event.keyCode;
        if (key === 32) {
            event.preventDefault();
        }
    })
})

// FUNCTION TO GET THE LIST OF REQUIRED PARAMETERS FOR VARIOUS SERVICE CHECKS
async function getServiceCheckParams(serviceCheckType) {

    const url = `/servicecheck/${serviceCheckType}/params`;
    const results = await fetch(url, {
        credentials: "same-origin",
        method: 'GET'
    })
    return await results.json();
}