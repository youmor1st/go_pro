$(document).ready(function() {
    $('#login-form').submit(function(event) {
        event.preventDefault();

        var username = $('#username').val();
        var password = $('#password').val();

        $.ajax({
            type: 'POST',
            url: '/login',
            contentType: 'application/json',
            data: JSON.stringify({
                username: username,
                password: password
            }),
            success: function(response) {
                $('#message').text('Login successful!');
                window.location.href = '/profile';
            },
            error: function(xhr, status, error) {
                $('#message').text(xhr.responseText);
            }
        });
    });
});
