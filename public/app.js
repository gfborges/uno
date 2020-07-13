new Vue({
    el:"#app",
    data :{
        name: "",
        wsGame: null,
        wsStack: null,  
        joined: false,
        canPass: false,
        myTurn: false,
        handCards: "",
        turnText:"",
        stack:{color:"", symbol:"", draws:0}
    },

    created : function(){
        let self = this;
        this.wsStack = new WebSocket('ws://' + window.location.host + '/ws/stack');
        this.wsGame = new WebSocket('ws://' + window.location.host + '/ws/game');
        this.wsGame.addEventListener("message", function(e){
            let msg = JSON.parse(e.data);
            console.log(msg.card.color, msg.card.symbol);

            let btn  = document.createElement("button");
            btn.id = msg.card.id;
            btn.className = `uno-card ${msg.card.color}`;
            btn.addEventListener("click", self.playCard);
            btn.innerHTML = `<div class="uno-inner"><div class="uno-frame"><span>${msg.card.symbol}</span></div></div>`;

            document.getElementById("hand").appendChild(btn);

        });
        this.wsStack.addEventListener("message", function(e) {
            let msg = JSON.parse(e.data);
            if(msg.name == "*"){
                self.turnText = msg.content
                return
            }
            let card  = msg.card;
            console.log("stack\n", msg)
            self.stack.color = card.color;
            self.stack.symbol = card.symbol;
            self.stack.draws = msg.draws
            let stackCard = document.getElementById("stack-card");
            let stackText = document.getElementById("stack-text");
            stackCard.className = `uno-card ${card.color}`;
            stackText.textContent = card.symbol;
            self.turnText = msg.name + " turn"
            if(!self.myTurn && msg.myTurn) {
                self.playSound("turn.mp3")
                Materialize.toast("It is your turn", 1000)
                setTimeout(function(){self.lostTurn()}, 45 * 1000)
            }
            self.myTurn = msg.myTurn
        });
        window.onbeforeunload = function (event) {
            self.wsGame.close()
            self.wsStack.close()
        }
    },

    methods: {
        join: function() {
            if(!this.name.match("^[a-z A-z0-9]+$")){
                Materialize.toast("invalid name", 1000);
                return;
            }
            console.log("joining")
            let txt = $('<p>').html(this.name).text()
            console.log(txt);
            let hello = JSON.stringify({
                name: txt,
            });
            this.wsStack.send(hello);
            this.wsGame.send(hello);
            this.joined = true;
        },

        isValidCard: function(card) {
            if(this.stack.draws > 0 && !card.symbol.startsWith("+")){
                return false
            }
            if(card.color == this.stack.color){
                return true;
            }
            if(card.symbol == "+4" || card.symbol == ""){
                return true;
            }
            if(this.stack.symbol == "+4" || this.stack.symbol == ""){
                return true;
            }
            if(card.symbol == this.stack.symbol){
                return true;
            }
            return false;
        },

        playCard: function(e) {
            if(!this.myTurn){
                Materialize.toast("Not your yout turn", 1000);
                return;
            }
            this.myTurn = false;
            let card = e.target;
            while(!card.className.includes("uno-card")){
                card = card.parentElement;
            }
            let cardObj = {
                name:this.name,
                card:{id:Number(card.id), 
                    color:card.classList[1], 
                    symbol:card.textContent }
            }
            if( !this.isValidCard(cardObj.card)){
                Materialize.toast("Invalid play", 1000);
                this.myTurn = true
                return;
            }
            this.playSound("sucess.mp3");
            this.canPass = false;
            if(card.className.includes("Wild")) {
                let colors = ["Blue", "Red", "Yellow", "Green"];
                let color = prompt("Choose a color\n1- Blue\n2- Red\n3- Yellow\n4- Green", "1");
                if(!color.match("^[1-4]$")){
                    this.myTurn = true
                    return
                }
                cardObj.card.color = colors[Number(color)-1]

            }
            card.remove();
            this.wsGame.send(JSON.stringify(cardObj)); 
        },  

        drawCard: function(e) {
            if(!this.myTurn || this.canPass){
                Materialize.toast("Cant draw card", 1000);
                return;
            }
            this.playSound("sucess.mp3");
            let msg = JSON.stringify({
                name:this.name,
                card:{id:-1, color:"*"},
            });
            this.wsGame.send(msg);
            this.canPass = true
        },

        playSound: function(sound){
            let audio = new Audio(sound);
            audio.play()
        },

        passTurn: function() {
            if(!this.canPass){
                Materialize.toast("Cant pass turn", 1000);
                return;
            }
            this.canPass = false;
            let msg = JSON.stringify({
                name:this.name,
                card:{id:-2, color:"*"},
            });
            this.wsGame.send(msg);
        }
    }
});