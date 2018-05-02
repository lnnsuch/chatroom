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
    self.userName = $('#group-1-1');
    self.sendButton.click(function () {
        self.sendMessage();
    });
    self.userName.on('click', '.other .time b', function () {
        self.clickName($(this).text());
    });
    self.roomTable.on('click', ".chatroom-tribe", function () {
        self.checkTab($(this));
    });
    let ws = new WebSocket("ws://" + document.location.host + "/ws");
    ws.onopen = function () {};
    ws.onmessage = function (evt) {
        self.takeMessage(JSON.parse(evt.data))
    };
    ws.onclose = function (evt) {};
    ws.onerror = function (evt) {};
    self.ws = ws;
};

chatRoot.prototype.takeMessage = function (data) {
    let text = '';
    if (data.status === 0) {
        text = this.takeTextMessage(data);
    } else if(data.status === -2) {
        text = this.takeSysMessage(data);
    }
    let idName = data.type + '-' + data.id;
    let current = $('#group-' + idName);
    if (current.length  === 0) {
        this.addWindow(data.type,data.id,data.name);
        current = $('#group-' + idName);
    }
    let group = $('#table-' + idName + ".current");
    if (group.length === 0) {
        group = $('#table-' + idName);
        let count = group.find('.count');
        let n = parseInt(count.text()) +1;
        if (n > 99) {
            count.text("99+");
        } else {
            count.text(n);
        }
        count.css('visibility', 'visible');
    }
    current.append(text);
    let scrollTop = current[0].scrollHeight;
    this.roomContent.scrollTop(scrollTop);
};

chatRoot.prototype.addWindow =function (type, id, name, isShow = false) {
    let c = isShow === true ? "current" : '';
    this.roomContent.append(`<div class="chatroom-item ${c}" id="group-${type}-${id}">
            </div>`);
    this.roomTable.append(`<li class="chatroom-tribe ${c}" id="table-${type}-${id}" data-id="${id}" data-type="${type}">
            <span class="name">${name}</span>
            <span class="count">0</span>
        </li>`);
};

chatRoot.prototype.takeTextMessage = function (data) {
    let c = this.name === data.name ? 'myself' : 'other';
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

chatRoot.prototype.clickName = function (name) {
  let id = window.localStorage.getItem(name);
  let self = this;
  if (!id) {
      $.get("/user/get", {name: name}).done(function (response) {
          if (response.status === 0) {
              window.localStorage.setItem(name, response.info);
              self.addTable(response.info, name);
          }
      })
  } else {
      self.addTable(id, name);
  }
};

chatRoot.prototype.addTable = function (id, name) {
    this.roomTable.find('.current').removeClass("current");
    this.roomContent.find('.current').removeClass("current");
    this.addWindow(2, id, name, true);
};

chatRoot.prototype.checkTab = function (self) {
    let id = self.data('id');
    let type = self.data('type');
    this.roomTable.find('.current').removeClass('current');
    self.find('.count').text('0');
    self.find('.count').css('visibility', 'hidden');
    self.addClass('current');
    this.roomContent.find('.current').removeClass('current');
    $('#group-' + type + '-' + id).addClass('current');
};

new chatRoot();