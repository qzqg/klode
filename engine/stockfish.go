package engine

import (
    "fmt"
    "os/exec"
    "strings"
)

// GetBestMove ترجع أفضل نقلة من Stockfish باستخدام سطر الأوامر
func GetBestMove(fen string) (string, error) {
    // مسار ملف Stockfish
    stockfishPath := "./stockfish/stockfish.exe"

    // 1. تشغيل Stockfish كعملية منفصلة
    cmd := exec.Command(stockfishPath)
    stdin, err := cmd.StdinPipe()
    if err != nil {
        return "", fmt.Errorf("فشل في إنشاء stdin: %v", err)
    }
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return "", fmt.Errorf("فشل في إنشاء stdout: %v", err)
    }

    // 2. بدء العملية
    if err := cmd.Start(); err != nil {
        return "", fmt.Errorf("فشل في تشغيل Stockfish: %v", err)
    }
    defer cmd.Process.Kill()

    // 3. إرسال أوامر UCI إلى Stockfish
    commands := []string{
        "uci",
        "setoption name UCI_LimitStrength value false",
        "setoption name UCI_Elo value 2800",
        "ucinewgame",
        "position fen " + fen,
        "go depth 18",
        "quit",
    }

    for _, command := range commands {
        if _, err := stdin.Write([]byte(command + "\n")); err != nil {
            return "", fmt.Errorf("فشل في إرسال الأمر '%s': %v", command, err)
        }
    }
    stdin.Close()

    // 4. قراءة مخرجات Stockfish
    output := make([]byte, 4096)
    n, err := stdout.Read(output)
    if err != nil {
        return "", fmt.Errorf("فشل في قراءة المخرجات: %v", err)
    }

    // 5. استخراج أفضل نقلة من المخرجات
    outputStr := string(output[:n])
    lines := strings.Split(outputStr, "\n")

    for _, line := range lines {
        if strings.HasPrefix(line, "bestmove") {
            parts := strings.Fields(line)
            if len(parts) >= 2 {
                bestMove := parts[1]
                // تجاهل "ponder" إذا وجد
                if bestMove != "(none)" {
                    return bestMove, nil
                }
            }
        }
    }

    return "", fmt.Errorf("لم يتم العثور على نقلة في المخرجات")
}