async function createUser() {
    console.log("Creating user!")
    const username = document.getElementById('username').value

    // GET BASE DATA
    const fd = new FormData();
    fd.append('username', username);
    fd.append('password', document.getElementById('password').value);
    fd.append('confirmPassword', document.getElementById('confirmPassword').value);
    fd.append('team', document.getElementById('team').value);
    fd.append('isAdmin', document.getElementById('isAdmin').value);

    // ATTEMPT TO CREATE THE SERVICE
    const response = await fetch(`/user/${username}`, {
        method: 'PUT',
        body: fd,
        credentials: 'same-origin'
    })

    window.location.href = '/admin/users/add';

}