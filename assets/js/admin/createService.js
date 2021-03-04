window.addEventListener('load', async () => {

    // PREVENT WHITESPACE FROM BEING INSERTED INTO THE SERVICE NAME FIELD
    const nameInput = document.getElementById('serviceName');
    nameInput.addEventListener('keypress', (event) => {
        const key = event.keyCode;
        if (key === 32) {
            event.preventDefault();
        }
    })

    // ANY TIME A DIFFERENT SERVICE CHECK TYPE IS CONFIGURED, UPDATE THE LIST OF REQUIRED ARGUMENTS
    const serviceCheckTypeInput = document.getElementById('serviceCheckType');
    serviceCheckTypeInput.addEventListener('change', async () => {
        const currentServiceCheckType = document.getElementById('serviceCheckType').value;

        // GET THE LIST OF ARGUMENTS FROM THE SERVER'S REST API
        const serviceCheckArguments = await getServiceCheckParams(currentServiceCheckType)

        // GET THE TARGET
        setParams(serviceCheckArguments);

    });

    // TRIGGER THE CHANGE EVENT ON THE SELECT
    const e = new Event('change');
    document.getElementById('serviceCheckType').dispatchEvent(e);
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

// CREATE A NEW FORM GROUP AND LABEL AND STUFF BASED ON THE SERVICE NAME
function buildNewFormGroup(paramName) {
    const fg = document.createElement('div');
    fg.classList.add('form-group', 'row');
    const col1 = document.createElement('div');
    col1.classList.add('col-md-3');
    const col2 = document.createElement('div');
    col2.classList.add('col-md-9');
    const label = document.createElement('label');
    label.htmlFor = `param-${paramName}`;
    label.innerText = paramName;
    const input = document.createElement('input');
    input.classList.add('form-control');
    input.type = 'text';
    input.id = `param-${paramName}`;

    col1.appendChild(label);
    col2.appendChild(input);
    fg.appendChild(col1);
    fg.appendChild(col2);
    return fg;

}

function clearParamFormGroups() {
    const paramContainer = document.getElementById('param-container');
    paramContainer.innerHTML = '';
}

function setParams(params) {

    console.log(`setting params ${params}`)

    // CLEAR EXISTING PARAMS
    clearParamFormGroups();

    // STORE NEW FORM GROUPS
    const formGroups = []
    for(const param of params){
        formGroups.push(buildNewFormGroup(param))
    }

    // ADD THEM
    const paramContainer = document.getElementById('param-container');
    for (const formGroup of formGroups) {
        paramContainer.append(formGroup);
    }

}

async function createService() {
    console.log("Creating service!")
    const serviceName = document.getElementById('serviceName').value

    // GET BASE DATA
    const fd = new FormData();
    fd.append('name', serviceName);
    fd.append('host', document.getElementById('serviceIP').value);
    fd.append('port', document.getElementById('servicePort').value);
    fd.append('proto', document.getElementById('serviceTransportProto').value);
    fd.append('checkType', document.getElementById('serviceCheckType').value);
    fd.append('username', document.getElementById('serviceUsername').value);
    fd.append('password', document.getElementById('servicePassword').value);

    // GET THE PARAM INFO
    const params = await getServiceCheckParams(document.getElementById('serviceCheckType').value);
    for (const param of params) {
        fd.append(param, document.getElementById(`param-${param}`).value);
    }
    console.log(fd);

    // ATTEMPT TO CREATE THE SERVICE
    const response = await fetch(`/service/${serviceName}`, {
        method: 'PUT',
        body: fd,
        credentials: 'same-origin'
    })

    window.location.href = '/admin/services/add';

}