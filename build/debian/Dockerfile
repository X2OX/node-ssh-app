FROM debian:latest

RUN apt -y update && apt -y upgrade && DEBIAN_FRONTEND=noninteractive apt -y install openssh-server sudo curl git nano wget zsh

# SSH
RUN echo 'root:password' | chpasswd
RUN echo "Port 22" >> /etc/ssh/sshd_config
RUN echo "PasswordAuthentication yes" >> /etc/ssh/sshd_config
RUN echo "PermitRootLogin yes" >> /etc/ssh/sshd_config

# neofetch
RUN wget -O neofetch https://github.com/dylanaraps/neofetch/raw/master/neofetch && chmod +x neofetch && mv neofetch /usr/bin

# oh-my-zsh
RUN chsh -s $(grep /zsh$ /etc/shells | tail -1) && sh -c "$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
RUN git clone --depth=1 https://github.com/zsh-users/zsh-autosuggestions ~/.oh-my-zsh/plugins/zsh-autosuggestions
RUN git clone --depth=1 https://github.com/zsh-users/zsh-syntax-highlighting ~/.oh-my-zsh/plugins/zsh-syntax-highlighting
RUN git clone --depth=1 https://github.com/zsh-users/zsh-completions ~/.oh-my-zsh/plugins/zsh-completions
RUN git clone --depth=1 https://github.com/zsh-users/zsh-history-substring-search ~/.oh-my-zsh/plugins/zsh-history-substring-search
RUN /bin/zsh -i -c 'omz update'

# config oh-my-zsh
RUN sed -i 's/ZSH_THEME="robbyrussell"/ZSH_THEME="kardan"/' ~/.zshrc
RUN echo "neofetch" >> ~/.zshrc
RUN /bin/zsh -c 'source ~/.zshrc'

# cloudflared
RUN arch=$(arch | sed s/aarch64/arm64/ | sed s/x86_64/amd64/) && curl -L --output cloudflared.deb "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-${arch}.deb"
RUN dpkg -i cloudflared.deb

COPY build/debian/run.sh run.sh
RUN chmod +x run.sh

CMD run.sh
