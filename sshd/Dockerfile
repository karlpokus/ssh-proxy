# Use an official Ubuntu as a base image
FROM ubuntu:20.04

# Prevent prompts during the package installation
ENV DEBIAN_FRONTEND=noninteractive

# Install SSH server
RUN apt-get update && \
    apt-get install -y openssh-server && \
    mkdir /var/run/sshd

# Set root password
RUN echo 'root:root' | chpasswd

# Allow root login and password authentication
RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
RUN sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config

# Expose SSH port
EXPOSE 22

# Start sshd
CMD ["/usr/sbin/sshd", "-D"]
