import net from 'net';
import fs from 'fs';
import cron from 'node-cron';

let host = '127.0.0.1';
let port = 6379;
const data = fs.readFileSync('json_data/data.json');
const jsonData = JSON.parse(data);

let socket = new net.Socket();


function scheduler_run() {
  // console.log("HI889")
  // cron.schedule('1 * * * *', () => {  
    // console.log("H2I")  // '8 0 * * *' => Run function every day at 8AM
    setInterval(ping_write, 10000)
    // console.log("H3I") 
  // });
}

function ping_write() {
  socket_writes(["PING"])
}

function socket_writes(data) {
  socket.write(encode(data));
  if (data[0] == "get") {
    return data[1]
  }
}


  socket.connect(port, host, () => {
    // console.log("1*\r\n$4\r\nping\r\n");
    // socket.write("*1\r\n$4\r\nping\r\n");
    console.log("connected")
    scheduler_run()
    // socket.write(encode(["ping"]));
  });
  socket.on('data', (data) => {
    let decoded_data = decode(data)
    console.log(decoded_data)
  
    jsonData.users.push({
      data: decoded_data,
    });
    fs.writeFileSync('json_data/data.json', JSON.stringify(jsonData));
      // console.log(`${data}`);
      socket.destroy();
      console.log(decoded_data)
      console.log(decoded_data == "PONG")
      console.log(decoded_data === "PONG")
      if(decoded_data === "PONG") {
          socket.destroy();
          clearInterval()
          console.log("I RAN");
          console.log(socket.closed);
          throw "stop execution";
      }
  });

function encode(data, simple_str=false) {
  if (typeof data == "undefined") {
    return `-{${String(data)}}\r\n`;
  }
  else if (typeof data == "string" && simple_str) {
    return `+${data}\r\n`;
  }
else if (typeof data == "number") {
  return `:${data}\r\n`;
}
else if (typeof data == "string") {
  return `$${data.length}\r\n${data}\r\n`
}
else if(data.constructor == Array) {
  let enc = `*${data.length}\r\n`
  for(let itm of data) {
    enc += encode(itm);
  }
  return enc;
  } else {
    console.log("Unknown type " + typeof data)
  }
}

function decode(data) {
  let marker = data.toString()[0]
  if (marker == "-") {
    return(data.toString().slice(1, -2))
  } else if(marker == "+") {
    return(data.toString().slice(1, -2))
  } else if(marker == ":") {
    return(parseInt(data.toString().slice(1, -2)))
  } else if(marker == "$") {
    let parts = data.toString().slice(1).split('\r\n');
    let ln = parseInt(parts[0])
    if (ln == -1 || ln == 0) {
      return ''
    } else {
      return parts[1].toString()
    }
  } else if(marker == "*") {
    let parts = data.toString().slice(1).split('\r\n');
    let ln = parseInt(parts[0])
    items = []
    for (let index = 1; index < ln + 1; index++) {
      items[index] = decode(`${parts[i]}\r\n`)
    }
    return items
  } else {

  }
}