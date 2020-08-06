Vue.component("conversation", {
	data: function () {
		return {}
	},
	props: ["index", "conversation", "isChoosed"],
	template: `
		<div class="conversation" v-bind:class="{ conversationClicked: this.isChoosed }">
			{{ conversation.lastMessage }}
		</div>`,
	methods: {
		select: function() {
			this.isChoosed = true
			this.conversation.lastMessage = "asd"
		},
	}
})

Vue.component("user-panel", {
	data: function () {
		return {}
	},
	template: `
		<div class="user-panel">
		</div>`, 
})

Vue.component("conversations", {
	data: function () {
		return {
			isChoosed: new Array(this.conversations.length).fill(false),
			choosedConversation: -1
		}
	},
	props: ["conversations"],
	template: `
		<div class="conversations">
			<div v-for="(conversation, index) in conversations">
				<conversation v-bind:isChoosed="isChoosed[index]" v-bind:conversation="conversation" v-bind:index="index" v-on:click.native="change(index)"/>
			</div>
		</div>`,
	methods: {
		change: function(index) {
			console.log(this.isChoosed)
			if (this.choosedConversation != index) {
				this.isChoosed = new Array(this.conversations.length).fill(false);
				this.isChoosed[index] = true;
			}
		}
	}
})

Vue.component("chat", {
	data: function() {
		return {}
	},
	props: ["messages"],
	template: `
		<div class="chat">
			<ul>
				<li v-for="message in messages">
					{{ message }}
				</li>
			</ul>
			<div id="control">
				<input v-model="currentMessage" v-on:keyup.enter="sendMessage"/>
				<button v-on:click="sendMessage" value="Send">Send</button>
			</div>
		</div>`,
	methods: {
		sendMessage: function() {
			// if (this.currentMessage !== "") {
			// 	this.messages.push(this.currentMessage);
			// 	this.currentMessage = "";
			// }
		},
	}
})

new Vue({
	el: '#app',
	template: `
	<div>
		<user-panel />
		<conversations v-bind:conversations="conversations"/>
		<!-- <chat /> -->
	</div>`,
	data: {
		conversations: [
			{lastMessage: "Добрый день, ок"},
			{lastMessage: "Привет! Во сколько ты приедешь?"},
			{lastMessage: "Отвечаю на вопросы которые задает сегодня..."},
			{lastMessage: "kek"},
			{lastMessage: "Поздравляем, Ваш номер подтвержден!"},
			{lastMessage: "Any other questions? Something didn`t..."},
			{lastMessage: "Mac Miller - Dunno"},
			{lastMessage: "Отвечаю на вопросы которые задает сегодня..."},
		],
		currentMessage: "",
		messages: [],
		connection: null,
		},
	created: function() {
		this.connection = new WebSocket("ws://"+window.location.host+"/ws");
	},
});
