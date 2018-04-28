let chatRoot = function () {
    this.documentInit();
    this.init();
};

chatRoot.prototype.documentInit = function () {
    // 禁止右键和刷新
    document.oncontextmenu = function () {
        return false;
    };
    document.onkeydown = function (event) {
        let e = event || window.event || arguments.callee.caller.arguments[0];
        if (e && e.keyCode == 116) {
            return false;
        }
    };
};
chatRoot.prototype.init = function () {
    let self = this;
    self.name = window.localStorage.getItem('name');
    self.roomContent = $('#roomContent');
    self.sendButton = $('#send');
    self.roomTable = $('#roomTable');
    self.sendButton.click(function () {
        self.sendMessage();
    });
    let ws = new WebSocket("ws://" + document.location.host + "/ws");
    ws.onopen = function () {};
    ws.onmessage = function (evt) {
        console.log(evt);
        // var node = document.getElementById('content');
        // var p = document.createElement('p');
        // p.innerHTML = evt.data;
        // node.appendChild(p);
        // node.scrollTop = node.scrollHeight - node.offsetHeight;
        self.takeMessage(JSON.parse(evt.data))
    };
    ws.onclose = function (evt) {};
    ws.onerror = function (evt) {};
    self.ws = ws;
};

chatRoot.prototype.takeMessage = function (data) {
    let current = this.roomContent.children('.current');
    let text = '';
    if (data.status === 0) {
        text = this.takeTextMessage(data);
    } else if(data.status === -2) {
        text = this.takeSysMessage(data);
    }
    current.append(text);
    let scrollTop = current[0].scrollHeight;
    this.roomContent.scrollTop(scrollTop);
};

chatRoot.prototype.takeTextMessage = function (data) {
    let c = this.name === data.name ? 'myself' : '';
    return `<div class="chatroom-log ${c}">
        <span class="avatar"><img src="https://avatars0.githubusercontent.com/u/30884897?s=40&v=4" alt="${data.name}"></span>
        <span class="time"><b data-id="Q-2xC-3e2q46">${data.name}</b>  ${new Date().toLocaleString()}</span>
        <span class="detail">${data.info}</span>
     </div>`
};

chatRoot.prototype.takeSysMessage = function (data) {
    return `<div class="chatroom-log system_log">
        <span><b>${data.info}</b>  ${new Date().toLocaleString()}</span>
     </div>`
};

chatRoot.prototype.sendMessage = function () {
    let t = this.sendButton.prev('textarea');
    let val = t.val();
    t.val('');
    let table = this.roomTable.children('.current');
    let info = {type: table.data('type'), id: table.data('id'), info: val};
    this.ws.send(JSON.stringify(info))
};


new chatRoot();