
// FUNCTION TO DELETE A SERVICE
async function deleteService(serviceName) {
    let url = `/blackteam/service/${serviceName}`
    let response = await fetch(url, {
        method: 'DELETE',
        credentials: 'same-origin'
    });
    window.location.href = '/blackteam/services/configure';
}


// FUNCTION TO UPDATE A SERICE
async function updateService(serviceName) {
    const url = `/blackteam/service/${serviceName}`;

    // BUILD A FORM DATA OBJECT
    const fd = new FormData();
    fd.append('host', document.getElementById(`${serviceName}-host`).value);
    fd.append('port',document.getElementById(`${serviceName}-port`).value);
    fd.append('transportProtocol', document.getElementById(`${serviceName}-transport-protocol`).value);
    fd.append('username', document.getElementById(`${serviceName}-username`).value);
    fd.append('password', document.getElementById(`${serviceName}-password`).value);

    console.log(fd)

    let response = await fetch (url, {
        method: 'PATCH',
        body: fd,
        credentials: 'same-origin'
    });

    window.location.href = '/blackteam/services/configure';

}