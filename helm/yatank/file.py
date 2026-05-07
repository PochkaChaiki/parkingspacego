# generate_ammo.py
def write_ammo(f, method, path, body=None, tag=""):
    # 1. Строка запроса
    req = f"{method} {path} HTTP/1.1\r\n"
    
    # 2. Заголовки (Phantom сам добавит общие заголовки из load.yaml)
    req += f"Host: parking-app-frontend.parking.svc.cluster.local:8080\r\n"
    req += f"Content-Type: application/json\r\n"
    req += f"Accept: application/json\r\n"
    req += f"Connection: keep-alive\r\n"
    
    # 3. Обработка тела POST-запроса
    if body is not None:
        body_bytes = body.encode('utf-8')
        req += f"Content-Length: {len(body_bytes)}\r\n"
        req += "\r\n"  # 🔑 ОБЯЗАТЕЛЬНЫЙ разделитель: конец заголовков
        req += body    # Тело запроса
    else:
        # Для GET указываем 0 и закрываем заголовки
        req += "Content-Length: 0\r\n"
        req += "\r\n"
        
    # 4. Точный подсчёт размера в байтах (UTF-8)
    size = len(req.encode('utf-8'))
    
    # 5. Запись в бинарном режиме (чтобы ОС не меняла \n на \r\n)
    f.write(f"{size} {tag}\n".encode('utf-8'))
    f.write(req.encode('utf-8'))

# Тела запросов (кириллица заменена на латиницу для теста, если нужно)
body1 = '{"client_name":"Test User","phone_number":"+79001234567","license_plate":"A123BC777","spot_number":42,"duration":"1h"}'
body2 = '{"client_name":"Load Tester","phone_number":"+79007654321","license_plate":"X999YY199","spot_number":15,"duration":"30m"}'

with open("ammo.txt", "wb") as f:
    write_ammo(f, "GET", "/api/sessions/+79001234567", tag="get_session_ok")
    write_ammo(f, "GET", "/api/sessions/+79999999999", tag="get_session_404")
    write_ammo(f, "POST", "/api/sessions", body=body1, tag="create_session")
    write_ammo(f, "POST", "/api/sessions", body=body2, tag="create_session_spot2")

print("✅ ammo.txt создан. Проверьте наличие пустых строк между заголовками и телом.")
