# MyDNS: Manual de Uso y Reemplazo de Dnsmasq

MyDNS ahora funciona como un DNS comodín (wildcard). Todo dominio que termine en `*.my.os` resolverá automáticamente a la IP de tu servidor (por defecto `192.168.1.114`). El proxy (MyProxy) se encarga del enrutamiento final al contenedor usando el header HTTP `Host`.

## 🚀 Cómo usar MyDNS

1.  **Construir imagen:**
    ```bash
    docker build -t mydns .
    ```
2.  **Ejecutar contenedor:**
    ```bash
    docker run -d \
        --name mydns \
        -p 53:53/udp -p 53:53/tcp \
        -e TARGET_IP=192.168.1.114 \
        --restart unless-stopped \
        mydns
    ```

## 🔄 Reemplazando Dnsmasq

Actualmente usas Dnsmasq con:
`address=/.my.os/192.168.1.114`

Para reemplazarlo totalmente:

1.  **Detén Dnsmasq en tu servidor Mint:**
    ```bash
    sudo systemctl stop dnsmasq
    sudo systemctl disable dnsmasq
    ```
2.  **Inicia MyDNS:** (ver paso anterior).
3.  **Configura tu Router/Clientes:**
    Apunta el DNS primario a la IP de tu servidor Mint (`192.168.1.114`).

## 🛠️ Flujo Simplificado

1.  **Resolución DNS**: El cliente pide `app.my.os`. MyDNS siempre responde `192.168.1.114`.
2.  **Tráfico al Proxy**: El navegador se conecta a `192.168.1.114:80`.
3.  **Enrutamiento**: MyProxy recibe la petición, lee que el `Host` es `app.my.os` y envía el tráfico al contenedor Docker correcto.
