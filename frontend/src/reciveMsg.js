"use strict";


import {
  ReciveFormatForJs,
} from "../wailsjs/go/main/App.js";

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
/**
 * Retrieves the current chat room id from the 'body' element's 'data-current-chat-room-id' attribute.
 *
 * @returns {string} The current chat room id.
 */
function getChatRoomId() {
  return document.getElementById("body").attributes["data-current-chat-room-id"]
    .value;
}


// Why the fuck do you stay stuck on fucking pending and dont resolve don't be a fucking special boy just be like eeveryone else like me. I can't recive the message from it that is a big problem, because I have to be done in 11 hours from now.
// FIXME: The comment above explains it well.
(function receive() {
  const channelId = getChatRoomId();
  console.info("receive");
  const msgPromise = ReciveFormatForJs(channelId)
  console.log(msgPromise);
  msgPromise.then((msg) => {
    console.log(msg);
  })
  setTimeout(receive, 1000);
})();
