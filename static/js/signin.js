var app = new Vue({
    el: "#app",
    data: {
        login: "",
        password: "",
    },
    methods: {
        doAuth: function() {
            var xhr = new XMLHttpRequest();

            xhr.open("POST", "/api/v1/auth", true)

            xhr.send(JSON.stringify({
                "login": this.login,
                "password": this.password,
            }));

            xhr.onreadystatechange = function() {
                if (xhr.readyState != 4) return;

                if (xhr.responseText === "ok") {
                    window.location.replace("/");
                }
            }
        }
    }
})