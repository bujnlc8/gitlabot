![build status](https://github.com/bujnlc8/gitlabot/actions/workflows/gitlabot.yml/badge.svg)

# gitlabot

## 将gitlab的webhook消息转发到企业微信机器人

## 使用

*   Build docker

        docker build -t gitlabot:0.0.1 .

*   Just Run it

    ```
     docker run -e listenAddr=0.0.0.0:9000 -p 9000:9000 -d  --name gitlabot  --restart always gitlabot:0.0.1

    ```

*   在`gitlab > Settings > Integrations` 新增webhook处的Secret Token填入企业微信机器人的推送key, URL处填入`http://127.0.0.1:9000/`, 当然也可以转发到此处。
