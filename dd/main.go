package main

import (
    "fmt"
    "log"
    "time"

    "chess-bot-go/browser"
    "chess-bot-go/engine"
)

func main() {
    fmt.Println("🚀 تشغيل بوت الشطرنج...")

    // 1. إعداد المتصفح
    wd, err := browser.SetupBrowser()
    if err != nil {
        log.Fatalf("❌ فشل في إعداد المتصفح: %v", err)
    }
    defer wd.Quit()
    fmt.Println("✅ تم إعداد المتصفح")

    // 2. تسجيل الدخول إلى Chess.com (ضع بياناتك هنا)
    email := "magharizk@gmail.com"   // غير هذا
    password := "tHn3ShQib3c_2Bq"          // غير هذا
    if err := browser.LoginToChess(wd, email, password); err != nil {
        log.Fatalf("❌ فشل في تسجيل الدخول: %v", err)
    }
    fmt.Println("✅ تم تسجيل الدخول")

    // 3. فتح مباراة (ضع رابط المباراة هنا)
    gameURL := "https://www.chess.com/game/live/123456789" // غير هذا
    if err := wd.Get(gameURL); err != nil {
        log.Fatalf("❌ فشل في فتح المباراة: %v", err)
    }
    time.Sleep(5 * time.Second)
    fmt.Println("✅ تم فتح المباراة")

    // 4. الحلقة الرئيسية - تلعب تلقائياً
    fmt.Println("♟️ البوت جاهز للعب...")
    for {
        // قراءة وضعية الرقعة
        fen, err := browser.GetFENFromBoard(wd)
        if err != nil {
            log.Printf("⚠️ خطأ في قراءة FEN: %v", err)
            time.Sleep(2 * time.Second)
            continue
        }

        // حساب أفضل نقلة من Stockfish
        bestMove, err := engine.GetBestMove(fen)
        if err != nil {
            log.Printf("⚠️ خطأ في حساب النقلة: %v", err)
            time.Sleep(2 * time.Second)
            continue
        }

        fmt.Printf("🎯 أفضل نقلة: %s\n", bestMove)

        // تنفيذ النقلة على اللوحة
        if err := browser.ExecuteMove(wd, bestMove); err != nil {
            log.Printf("⚠️ خطأ في تنفيذ النقلة: %v", err)
        }

        // انتظار دور الخصم
        time.Sleep(3 * time.Second)
    }
}