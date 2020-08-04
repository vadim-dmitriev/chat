var app = new Vue({
    el: '#app',
    template: `
      <div>
        <div id="message-list">
          <ul>
            <li v-for="message in messages">
              {{ message }}
            </li>
          </ul>
        </div>
        <div id="control">
          <input v-model="currentMessage"
                v-on:keyup.enter="sendMessage"/>
          <button v-on:click="sendMessage" value="Send">Send</button>
        </div>
      </div>`,
    data: {
      currentMessage: "",
      messages: [],
      connection: null,
    },
    created: function() {
      this.connection = new WebSocket("ws://"+window.location.host+"/ws");
    },
    methods: {
      sendMessage: function() {
        if (this.currentMessage !== "") {
          this.messages.push(this.currentMessage);
          this.currentMessage = "";
        }
      },
    }

});
