Vue.component("conversation", {
	data: function () {
		return {
			// isChoosed: true
		}
	},
	props: ["conversation", "choosedConversation", "name"],
	template: `
		<div class="conversation" v-bind:class="{ conversationClicked: this.isChoosed }">
			<strong>{{ name }}</strong>
			<br/><br/>
			{{ lastMessage }}
			<br/><br/>
			{{ lastMessageTime }}
		</div>`,
	methods: {},
	computed: {
		isChoosed: function() {
			return this.name === this.choosedConversation
		},
		lastMessage: function() {
			return this.conversation.messages[0].value
		},
		lastMessageTime: function() {
			return this.conversation.messages[0].time
		}

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
			choosedConversation: ""
		}
	},
	props: ["conversations"],
	template: `
		<div class="conversations">
			<div v-for="conversation in Object.keys(conversations)">
				<conversation :name="conversation" :conversation="conversations[conversation]"  v-on:click.native="change(conversation)"
							  :choosedConversation="choosedConversation"
				/>
			</div>
		</div>`,
	methods: {
		change: function(convName) {
			if (this.choosedConversation !== convName) {
				this.choosedConversation = convName;
				this.$emit("change-conversation", convName);
			}
		}
	},
});

Vue.component("chat", {
	data: function() {
		return {
			currentMessage: "",
		}
	},
	props: ["conversation"],
	template: `
		<div class="chat" v-show="isActive">
			<ul>
				<li v-for="message in messages.slice().reverse()">
					{{ message.value }} - {{   message.time }}
				</li>
			</ul>
			<div>
				<input v-model="currentMessage" v-on:keyup.enter="sendMessage"/>
				<button v-on:click="sendMessage" value="Send">Send</button>
			</div>
		</div>
	`,
	methods: {
		sendMessage: function() {
			if (this.currentMessage !== "") {
				this.$emit("send-message", this.currentMessage)
				this.conversation.messages.unshift({
					value: this.currentMessage,
					time: new Date().toString(),
				})
				this.currentMessage = ""
			}
		},
	},
	computed: {
		messages: function() {
			return this.conversation.messages
		},
		isActive: function() {
			return (this.conversation !== {})
		}
	},
});

var app = new Vue({
	el: '#app',
	data: {
		conversations: {},
		currentConversation: {},
		currentConversationName: "",
		ws: null,
	},

	created: function() {
		t = this

		this.ws = new WebSocket("ws://"+window.location.host+"/api/v1/ws");

		this.ws.onopen = function() {
			t.getConversaions();
		}

		this.ws.onmessage = function(event) {
			message = JSON.parse(event.data)
			console.log("NEW MESSAGE FROM SERVER:", message);

			switch (message.action) {
			case "conversations":
				t.conversations = message.conversations;
				console.log(t.conversations)
				break;

			case "searchUser":
				if (message.isUserExists) {
					newConv = t.conversations
					newConv[message.newConversationWith] = {
						is_dialog: true,
						messages: [
							{
								value: "Нет сообщений...",
								sender: "",
								time: "",
							}
						]
					}
					t.conversations = newConv
					console.log(t.conversations)
					t.currentConversation = t.conversations[message.newConversationWith]
					t.currentConversationName = message.newConversationWith
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
					conversationName: this.currentConversationName,
					message: message
				})
			);
		},
		changeConversation: function(currentChoosedConversation) {
			this.currentConversation = this.conversations[currentChoosedConversation];
			this.currentConversationName = currentChoosedConversation
		}
	}
});
