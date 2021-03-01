
// FUNCTION TO DELETE A SERVICE
async function deleteService(serviceName) {
    let url = `/service/${serviceName}`
    let response = await fetch(url, {
        method: 'DELETE',
        credentials: 'same-origin'
    });
    window.location.href = '/admin/services/configure';
}