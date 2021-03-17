
// FUNCTION TO DELETE A SERVICE
async function deleteUser(username) {
    let url = `/blackteam/user/${username}`
    let response = await fetch(url, {
        method: 'DELETE',
        credentials: 'same-origin'
    });
    window.location.href = '/blackteam/users/configure';
}


// FUNCTION TO UPDATE A SERICE
async function updateService(username) {
    const url = `/blackteam/user/${username}`;

    // BUILD A FORM DATA OBJECT
    const fd = new FormData();

    // TODO FINISH THIS

    console.log(fd)

    let response = await fetch (url, {
        method: 'PATCH',
        body: fd,
        credentials: 'same-origin'
    });

    window.location.href = '/blackteam/users/configure';

}