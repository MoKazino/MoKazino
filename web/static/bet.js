function multiply(elmid, fix = 0) {
    elm = document.getElementById(elmid)
    elm.value = Number(elm.value * 2).toFixed(fix)
    domath()
}
function divide(elmid, fix = 0) {
    elm = document.getElementById(elmid)
    elm.value = Number(elm.value * 0.5).toFixed(fix)
    domath()
}

function setmaxwd() {
    bal = document.getElementById("balance").value
    document.getElementById("withdrawamount").value = lastbal
}

function add(elmid, fix = 0) {
    elm = document.getElementById(elmid)
    // No, I don't like JS.
    elm.value = Number(elm.value - (-5)).toFixed(2)
    domath()
}
function substract(elmid, fix = 0) {
    elm = document.getElementById(elmid)
    elm.value = Number(elm.value - 5).toFixed(2)
    domath()
}

var didinit = false
var lastbal = 0;
function init() {
    fetch("/api/v1/profile")
        .then(r => r.json())
        .then(r => {
            if (r.user.currency == "BTC") {
                document.getElementById("logo").src = "/img/bitcoin.svg"
            } else {
                document.getElementById("logo").src = "/img/monero.svg"
            }
//            if (!didinit) {
//                didinit = true;
//                for (let i = 0; i < r.bets.length; i++) {
//                    addBetToLog(r.win, r.nonce, Number(r.profit/1e12).toFixed(8), Number(chance).toFixed(2), Number(r.roll/100).toFixed(2))
//                }
//            }
            if (r.user.currency == "BTC") {
                document.getElementById("balance").value = (r.user.bitcoinbalance/1e8).toFixed(8)
                lastbal = r.user.bitcoinbalance/1e8
            } else {
                document.getElementById("balance").value = (r.user.balance/1e12).toFixed(8)
                lastbal = r.user.balance/1e12
            }
            //document.getElementById("apikey").innerText = r.user.cookie
            if (r.user.currency == "BTC") {
                document.getElementById("depoaddr").innerText = r.user.bitcoinaddress
            } else {
                document.getElementById("depoaddr").innerText = r.user.moneroaddress
            }
            document.getElementById("nonce").value = r.user.nonce
            document.getElementById("clientseed").value = r.user.clientseed
            document.getElementById("usernamechat").innerText = r.user.username
            document.getElementById("refcodeuserid").innerText = r.user.id
            if (r.user.currency == "BTC") {
                document.getElementById("investbalance").innerText = (r.user.investbalancebitcoin/1e8).toFixed(12)
            } else {
                document.getElementById("investbalance").innerText = (r.user.investbalance/1e12).toFixed(12)    
            }
            
            document.getElementById("cur1").innerText = r.user.currency
            document.getElementById("cur2").innerText = r.user.currency
            document.getElementById("cur3").innerText = r.user.currency
            document.getElementById("cur4").innerText = r.user.currency
            document.getElementById("cur5").innerText = r.user.currency
            // document.getElementById("cur6").innerText = r.user.currency
            document.getElementById("cur7").innerText = r.user.currency

            document.getElementById("majesticdepo"+r.user.currency.toLowerCase()).style.visibility = "hidden";
            fetch("/api/v1/serverstats")
            .then(r2 => r2.json())
            .then(r2 => {
                document.getElementById("statsmaxprofit").innerText = r2.maxprofit+" "+r.user.currency;
                document.getElementById("statsuserbalance").innerText = r2.userbalance+" "+r.user.currency;
                document.getElementById("statsserverbalance").innerText = r2.serverbalance+" "+r.user.currency;
            })
        })
        .catch(e => {
            console.log(e)
            deleteCookies()
            window.location.href = "/"
        })
}

function send() {
    fetch('/api/v1/chat/send', {
        method: 'POST',
        body: new URLSearchParams({
            'text': document.getElementById("chatmsg").value
        })
    })
    .then(r=> {
        document.getElementById("chatmsg").value = "";
    })
}

function watchForMessages() {
    fetch('/api/v1/chat/read')
    .then(r => r.text())
    .then(r => {
        let chatlog = document.getElementById("chatlog")
        // <li class="list-group-item">msg</li>
        let msg = document.createElement("li")
        msg.classList="list-group-item"
        msg.innerText = r
        chatlog.appendChild(msg)
        scrollToBottom('chatlogdiv')
        setTimeout(watchForMessages)
    })
}
setTimeout(watchForMessages)
function saveprovablyfair() {
    fetch('/api/v1/profile_update', {
        method: 'POST',
        body: new URLSearchParams({
            'clientseed': document.getElementById("clientseed").value,
        })
    }).then(r => {
        init()
    })
}
function scrollToBottom (id) {
    var div = document.getElementById(id);
    div.scrollTop = div.scrollHeight - div.clientHeight;
}
function withdraw() {
    fetch('/api/v1/withdraw', {
        method: 'POST',
        body: new URLSearchParams({
            'amount': Number(document.getElementById("withdrawamount").value*1e12).toFixed(0),
            'address': document.getElementById("withdrawaddress").value,
        })
    })
    .then(r => r.text())
    .then(r => {
        let a = document.createElement("div")
        a.classList = "alert alert-danger"
        a.id = "betstats"
        a.innerHTML = r
        document.getElementById("betstats").outerHTML = a.outerHTML
        init()
    })
}

function bet(betlo) {
    placebet(
        document.getElementById("chance").value,
        document.getElementById("amount").value,
        betlo
    )
}

function placebet(chance, amount, betlo) {
    fetch('/api/v1/bet', {
        method: 'POST',
        body: new URLSearchParams({
            'amount': Number(amount*1e12).toFixed(0),
            'chance': Number(chance*100).toFixed(0),
            'betlo': betlo,
        })
    })
    .then(r => r.text())
    .then(resp => {
        try {
            let r = JSON.parse(resp)
            document.getElementById("nonce").value = r.nonce
            document.getElementById("balance").value = (r.balance/1e12).toFixed(8)
            lastbal = r.balance/1e12;
            let a = document.createElement("div")
            let hm
            if (r.win) {
                hm = "success"
            } else {
                hm = "danger"
            }
            a.classList = "alert alert-"+hm
            a.id = "betstats"
            if (r.win) {
                a.innerHTML = "You have won <code>"+(r.profit/1e12).toFixed(8)+`</code> ${ document.getElementById("cur1").innerText }<br />${ Number(r.roll/100).toFixed(2) } Target: ${ Number(chance*100).toFixed(0) }`
            } else {
                a.innerHTML = "You have lost <code>"+(r.profit/1e12).toFixed(8)+`</code> ${ document.getElementById("cur1").innerText }<br />${ Number(r.roll/100).toFixed(2) } Target: ${ Number(chance*100).toFixed(0) }`
            }
            document.getElementById("betstats").outerHTML = a.outerHTML
            addBetToLog(r.win, r.nonce, Number(r.profit/1e12).toFixed(8), Number(chance).toFixed(2), Number(r.roll/100).toFixed(2))
        } catch(err) {
            console.log(err)
            let a = document.createElement("div")
            a.classList = "alert alert-danger"
            a.id = "betstats"
            a.innerHTML = resp
            document.getElementById("betstats").outerHTML = a.outerHTML
        }
    })
} 

function domath() {
    let bet = document.getElementById("amount").value
    let chance = document.getElementById("chance").value
    let multi = 1/(chance/100) * 0.99
    document.getElementById("multiplynum").innerText = Number(multi).toFixed(2)
    document.getElementById("profitwin").value = Number((bet*multi)-bet).toFixed(8)
}

function addBetToLog(win, id, profit, chance, roll) {
    let table = document.getElementById('bettable')
    table.children
    //<tr>
    //    <th scope="row">1</th>
    //    <td>Mark</td>
    //    <td>Otto</td>
    //    <td>@mdo</td>
    //
    let tr = document.createElement('tr')
    if (win) {
        tr.classList = "table-success"
    } else {
        tr.classList = "table-danger"
    }
    let th_betid = document.createElement('th')
    th_betid.attributes["scope"] = "row"
    th_betid.setAttribute("scope", "row")
    th_betid.innerText = id
    let th_profit = document.createElement('td')
    th_profit.innerText = profit
    let th_chance = document.createElement("td")
    th_chance.innerText = chance
    let th_roll = document.createElement("td")
    th_roll.innerText = roll
    tr.appendChild(th_betid)
    tr.appendChild(th_profit)
    tr.appendChild(th_chance)
    tr.appendChild(th_roll)
    while (table.children.length > 15) {
        table.removeChild(table.children[table.children.length-1])
    }
    table.insertBefore(tr, table.firstChild)
}
setTimeout(domath)

function deleteCookies() {
    var allCookies = document.cookie.split(';');
    
    // The "expire" attribute of every cookie is 
    // Set to "Thu, 01 Jan 1970 00:00:00 GMT"
    for (var i = 0; i < allCookies.length; i++)
        document.cookie = allCookies[i] + "=;expires="
        + new Date(0).toUTCString();
}

function invest() {
    fetch('/api/v1/invest', {
        method: 'POST',
        body: new URLSearchParams({
            'amount': Number(document.getElementById("investamount").value*1e12).toFixed(0),
            'address': document.getElementById("investamount").value,
        })
    })
    .then(r => r.text())
    .then(r => {
        let a = document.createElement("div")
        a.classList = "alert alert-danger"
        a.id = "investstatus"
        a.innerHTML = r
        document.getElementById("investstatus").outerHTML = a.outerHTML
        init()
    })
}

function closeinvest() {
    fetch('/api/v1/closeinvest', {})
    .then(r => r.text())
    .then(r => {
        let a = document.createElement("div")
        a.classList = "alert alert-danger"
        a.id = "investstatus"
        a.innerHTML = r
        document.getElementById("investstatus").outerHTML = a.outerHTML
        init()
    })
}

function majesticbank(currency) {
    fetch('/api/v1/external/majesticbank/exchange?receive_curency='+currency, {})
    .then(r => r.json())
    .then(r => {
        console.log(r)
        document.getElementById("majcurid").innerText = r.from_currency
        document.getElementById("majexchangeid").innerText = r.trx
        document.getElementById("majaddress").innerText = r.address
        let cur = "bitcoin";
        if (r.from_currency == "LTC") {
            cur = "litecoin";
        }
        document.getElementById("majqrcode").src = "/api/v1/depositqr/custom?qr="+cur+":"+r.address
        document.getElementById("majexchangerate").innerText = "1 "+r.from_currency+" = "+r.receive_amount+" "+r.receive_currency
        document.getElementById("majexpire").innerText = r.expiration
    })
}