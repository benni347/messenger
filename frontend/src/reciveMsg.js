"use strict";

import amqp from 'amqplib/callback_api';

export class App {
  constructor() {
    // Initialize properties
    this.rabbitMqHost = this.getRabbitMqHost();
    this.rabbitMqAdmin = this.getRabbitMqAdmin();
    this.rabbitMqPassword = this.getRabbitMqPassword();
  }

  getRabbitMqHost() {
    // Implement your method to retrieve RabbitMQ Host
  }

  getRabbitMqAdmin() {
    // Implement your method to retrieve RabbitMQ Admin
  }

  getRabbitMqPassword() {
    // Implement your method to retrieve RabbitMQ Password
  }

  receive(channelId) {
    const amqpUrl = `amqp://${this.rabbitMqAdmin}:${encodeURIComponent(this.rabbitMqPassword)}@${this.rabbitMqHost}:5672/`;

    amqp.connect(amqpUrl, (err, conn) => {
      if (err) {
        console.error("Failed to connect to RabbitMQ", err);
        return;
      }

      conn.createChannel((err, channel) => {
        if (err) {
          console.error("Failed to open a channel", err);
          return;
        }

        const queueName = channelId;
        channel.assertQueue(queueName, { durable: true }, (err) => {
          if (err) {
            console.error("Failed to declare a queue", err);
            return;
          }

          channel.consume(
            queueName,
            (msg) => {
              if (msg !== null) {
                console.log(msg.content.toString());
                channel.ack(msg);
              }
            },
            { noAck: false }
          );
        });
      });
    });
  }
}

// Usage
const app = new App();
app.receive("myChannelId");

