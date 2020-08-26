var app = new Vue({
    el: "#app",
    data: {
        login: "",
        password: "",
    },
    methods: {
        doRegister: function() {
            var xhr = new XMLHttpRequest();
            xhr.open("POST", "/api/v1/register", false);

            xhr.send(
                JSON.stringify({
                    "username": this.login,
                    "password": this.password,
                })
            );

            if (xhr.status == 200) {
                window.location.replace("/signin");
            }
        }
    }
});