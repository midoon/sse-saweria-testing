import express from "express";
import cors from "cors";
import bodyParser from "body-parser";
import { connect } from "amqplib";
import { v4 as uuidv4 } from "uuid";

const app = express();
app.use(bodyParser.urlencoded({ extended: false }));
app.use(bodyParser.json());
const corsConfig = {
  credentials: true,
  origin: true,
};
app.use(cors(corsConfig));
const port = 5000;

const sendMessageToRMQ = async (message) => {
  try {
    const connection = await connect("amqp://guest:guest@localhost:5672/");
    const channel = await connection.createChannel();

    channel.publish("sse", "saweria", Buffer.from(message), {
      headers: {
        ContentType: "text/plain",
      },
    });

    await channel.close();
    await connection.close();
  } catch (e) {
    console.log(e);
  }
};

app.post("/", async (req, res) => {
  const dataJson = {
    id: uuidv4(),
    name: req.body.donator_name,
    email: req.body.donator_email,
    amount: req.body.amount_raw,
    message: req.body.message,
    date: new Date(),
  };

  await sendMessageToRMQ(JSON.stringify(dataJson));
  return res.status(200);
});

app.listen(port, () => {
  console.log(`running in port ${5000}`);
});
