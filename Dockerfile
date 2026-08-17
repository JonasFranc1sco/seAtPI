FROM golang:1.24.6-bullseye

RUN apt update && apt install -y --no-install-recommends \
                    git \
                    zsh \
                    curl \
                    wget \
                    fonts-powerline \
                    lynx

RUN go env -w GO111MODULE=on          # on em vez de auto

RUN useradd -ms /bin/bash appuser
USER appuser

# GOPATH separado do projeto
ENV GOPATH=/home/appuser/go
WORKDIR /home/appuser/app

COPY ./.docker/cert.crt /certs/ca-certificates.crt
ENV SSL_CERT_FILE=/certs/ca-certificates.crt

RUN sh -c "$(wget -O- https://github.com/deluan/zsh-in-docker/releases/download/v1.1.5/zsh-in-docker.sh)" -- \
    -t https://github.com/romkatv/powerlevel10k \
    -p git \
    -p git-flow \
    -p https://github.com/zdharma-continuum/fast-syntax-highlighting \
    -p https://github.com/zsh-users/zsh-autosuggestions \
    -p https://github.com/zsh-users/zsh-completions \
    -a 'export TERM=xterm-256color'

RUN echo '[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh' >> ~/.zshrc && \
    echo 'HISTFILE=/home/appuser/zsh/.zsh_history' >> ~/.zshrc && \
    echo 'eval "$(pdm --pep582)"' >> ~/.zshrc && \
    echo 'eval "$(pdm --pep582)"' >> ~/.bashrc

CMD ["tail", "-f", "/dev/null"]