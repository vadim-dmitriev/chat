Vue.component("conversation", {
	data: function () {
		return {}
	},
	props: ["index", "conversation", "isChoosed", "name"],
	template: `
		<div class="conversation" v-bind:class="{ conversationClicked: this.isChoosed }">
			<strong>{{ name }}</strong>
			<br/><br/>
			{{ conversation.last_message.value }}
			<br/><br/>
			{{ conversation.last_message.time }}
		</div>`,
	methods: {
		select: function() {
			// this.isChoosed = true
			// this.conversation.lastMessage.value = "asd"
		},
	}
})

Vue.component("search-user-panel", {
	data: function () {
		return {
			username: ""
		}
	},
	props: ["isActive"],
	template: `
		<div v-show="isActive">
			<button v-on:click="hidePanel">x</button>
			<input v-model="username" />
			<button v-on:click="searchUser">Search</button>
		</div>
	`,
	methods: {
		hidePanel: function() {
			this.$emit("hidePanel");
		},
		searchUser: function() {
			if (this.username === "") {
				return
			}

			this.$emit("search-user", this.username)
			this.username = ""
			this.hidePanel()
		}
	}
})

Vue.component("user-panel", {
	data: function () {
		return {
			isUserSearchPanelActive: false,
		}
	},
	template: `
		<div class="user-panel">
			<div style="border: 1px black solid;" v-on:click="showSearchUserPanel">
				<p>Найти пользователя</p>
			</div>
			<search-user-panel v-bind:isActive="isUserSearchPanelActive"
							   v-on:hidePanel="hideSearchUserPanel"
							   v-on:search-user="searchUser"/>
		</div>`, 
	methods: {
		showSearchUserPanel: function() {
			this.isUserSearchPanelActive = true;
		},
		hideSearchUserPanel: function() {
			this.isUserSearchPanelActive = false;
		},
		searchUser: function(username) {
			this.$emit("search-user", username)
		}
	}
})

Vue.component("conversations", {
	data: function () {
		return {
			// isChoosed: new Array(this.conversations.length).fill(false),
			// choosedConversation: -1
		}
	},
	props: ["conversations"],
	template: `
		<div class="conversations">
			<div v-for="conversation in Object.keys(conversations)">
				<conversation :name="conversation" :conversation="conversations[conversation]"  v-on:click.native="change(index)"/>
			</div>
		</div>`,
	methods: {
		change: function(index) {
			// if (this.choosedConversation != index) {
			// 	this.isChoosed = new Array(this.conversations.length).fill(false);
			// 	this.isChoosed[index] = true;
			// 	this.$emit("change-conversation", this.conversations[index].name);
			// }
		}
	},
	watch: {
		// conversations: function() {
		// 	this.change(this.conversations.length-1)
		// }
	}
});

Vue.component("chat", {
	data: function() {
		return {
			currentMessage: "",
		}
	},
	props: ["messages"],
	template: `
		<div class="chat">
			<ul>
				<li v-for="message in messages">
					{{ message }}
				</li>
			</ul>
			<div>
				<input v-model="currentMessage" v-on:keyup.enter="sendMessage"/>
				<button v-on:click="sendMessage" value="Send">Send</button>
			</div>
		</div>`,
	methods: {
		sendMessage: function() {
			if (this.currentMessage !== "") {
				this.$emit("send-message", this.currentMessage)
				this.messages.push(this.currentMessage)
				this.currentMessage = ""
			}
		},
	}
});

var app = new Vue({
	el: '#app',
	data: {
		conversations: {},
		currentConversation: "",
		ws: null,
	},

	created: function() {
		t = this

		this.ws = new WebSocket("ws://"+window.location.host+"/api/v1/ws");

		this.ws.onopen = function(event) {
			t.getConversaions();
		}

		this.ws.onmessage = function(event) {
			message = JSON.parse(event.data)
			console.log(message);

			switch (message.action) {
			case "conversations":
				t.conversations = message.conversations;
				console.log(t.conversations)

				break;

			case "searchUser":
				if (message.isUserExists) {
					t.conversations.push({
						name: message.newConversationWith,
					});
				}
				break;

			case "newMessage":
				alert(message.value);
				break;
			}
		}

	},
	destroyed: function() {
		this.ws.close()
	},
	methods: {
		getConversaions: function () {
			this.ws.send(
				JSON.stringify({
					action: "getConversations"
				})
			);

		},
		searchUser: function(username) {
			this.ws.send(
				JSON.stringify({
					action: "searchUser",
					username: username
				})
			);
		},
		sendMessage: function(message) {
			this.ws.send(
				JSON.stringify({
					action: "sendMessage",
					conversationName: this.currentConversation,
					message: message
				})
			);
		},
		changeConversation: function(currentChoosedConversation) {
			this.currentConversation = currentChoosedConversation;
		}
	}
});
