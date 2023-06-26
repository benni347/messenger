"use strict";

import { ReciveFormatForJs } from "../wailsjs/go/main/App.js";

// export class App {
//   constructor() {
//     // Initialize properties
//     this.rabbitMqHost = this.getRabbitMqHost();
//     this.rabbitMqAdmin = this.getRabbitMqAdmin();
//     this.rabbitMqPassword = this.getRabbitMqPassword();
//   }
//
//   getRabbitMqHost() {
//     return GetRabbitMqHost();
//   }
//
//   getRabbitMqAdmin() {
//     return GetRabbitMqAdmin();
//   }
//
//   getRabbitMqPassword() {
//     return GetRabbitMqPassword();
//   }
//
//   receive(channelId) {
//     const amqpUrl = `amqp://${this.rabbitMqAdmin}:${encodeURIComponent(this.rabbitMqPassword)}@${this.rabbitMqHost}:5672/`;
//
//     amqp.connect(amqpUrl, (err, conn) => {
//       if (err) {
//         console.error("Failed to connect to RabbitMQ", err);
//         return;
//       }
//
//       conn.createChannel((err, channel) => {
//         if (err) {
//           console.error("Failed to open a channel", err);
//           return;
//         }
//
//         const queueName = channelId;
//         channel.assertQueue(queueName, { durable: true }, (err) => {
//           if (err) {
//             console.error("Failed to declare a queue", err);
//             return;
//           }
//
//           channel.consume(
//             queueName,
//             (msg) => {
//               if (msg !== null) {
//                 console.log(msg.content.toString());
//                 channel.ack(msg);
//               }
//             },
//             { noAck: false }
//           );
//         });
//       });
//     });
//   }
// }
//
// // Usage
// const app = new App();

//
// app.receive(getChatRoomId());
//
//
const socket = new WebSocket("ws://s.rabbitmq.cedric.5y5.one/websocket");

socket.onmessage = function(event) {
  const otherMessage = event.data;
  const messageLog = document.getElementById("message-log");
  const messageDiv = document.createElement("div");
  const messageUsernameDiv = document.createElement("div");
  const messageTextDiv = document.createElement("div");
  const otherUsername = "other";
  messageUsernameDiv.innerHTML = otherUsername;
  messageTextDiv.innerHTML = otherMessage;
  messageTextDiv.className = "text";
  messageUsernameDiv.className = "username";
  messageDiv.className = "message";
  messageDiv.appendChild(messageUsernameDiv);
  messageDiv.appendChild(messageTextDiv);
  messageLog.appendChild(messageDiv);
  messageLog.scrollTop = messageLog.scrollHeight;
  // let data = JSON.parse(event.data);
  // console.log(`Received message from ${data.queue}: ${data.message}`);
};

