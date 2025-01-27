# Use a PHP base image with Apache
FROM php:8.1-apache

# Install dependencies
RUN apt-get update && apt-get install -y \
    libpcre3-dev \
    zlib1g-dev \
    libpq-dev \
    zip \
    unzip \
    git && \
    docker-php-ext-install pdo pdo_pgsql sockets && \
    pecl install phalcon && \
    docker-php-ext-enable phalcon

# Enable Apache mod_rewrite
RUN a2enmod rewrite

# Add Authorization header forwarding
RUN echo 'SetEnvIf Authorization "(.*)" HTTP_AUTHORIZATION=$1' >> /etc/apache2/conf-enabled/auth.conf

# Set the working directory inside the container
WORKDIR /var/www/html

# Copy your application's code into the container
COPY . /var/www/html

# Install Composer
COPY --from=composer:2 /usr/bin/composer /usr/bin/composer

# Install PHP dependencies
RUN composer install --no-dev --optimize-autoloader

# Set permissions
RUN chown -R www-data:www-data /var/www/html && \
    chmod -R 755 /var/www/html

# Expose port 80
EXPOSE 80

# Start Apache server
CMD ["apache2-foreground"]
