package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// ─── CONFIG ───────────────────────────────────────────────────────────────────

const (
	WhatsAppTarget = "083899782135"
	FonnteAPIURL   = "https://api.fonnte.com/send"
	ServerName     = "Dewata Nation Roleplay"
)

func getFonnteToken() string {
	token := os.Getenv("FONNTE_TOKEN")
	if token == "" {
		token = "rAi9rzrezVBFFfe5w1Gp" // Ganti dengan token Fonnte kamu
	}
	return token
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return port
}

// ─── STRUCTS ──────────────────────────────────────────────────────────────────

type AdminForm struct {
	// OOC (Out of Character)
	OOCNamaAsli     string `json:"ooc_nama_asli"`
	OOCUmur         string `json:"ooc_umur"`
	OOCWilayah      string `json:"ooc_wilayah"`
	OOCDiscord      string `json:"ooc_discord"`
	OOCNoHP         string `json:"ooc_no_hp"`
	OOCSudahBerapa  string `json:"ooc_sudah_berapa"`

	// IC (In Character)
	ICNamaKarakter  string `json:"ic_nama_karakter"`
	ICUmurKarakter  string `json:"ic_umur_karakter"`
	ICPekerjaanIC   string `json:"ic_pekerjaan_ic"`
	ICLevelChar     string `json:"ic_level_char"`
	ICWarnedBanned  string `json:"ic_warned_banned"`

	// Pengalaman Admin
	PengalamanAdmin    string `json:"pengalaman_admin"`
	ServerSebelumnya   string `json:"server_sebelumnya"`
	LamaBermain        string `json:"lama_bermain"`
	KeahlianKhusus     string `json:"keahlian_khusus"`

	// Motivasi & Komitmen
	Motivasi           string `json:"motivasi"`
	KontribusiRencana  string `json:"kontribusi_rencana"`
	KetersediaanWaktu  string `json:"ketersediaan_waktu"`
	SkenarioHandler    string `json:"skenario_handler"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ─── HTML TEMPLATE ────────────────────────────────────────────────────────────

const htmlTemplate = `<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
  <title>Pendaftaran Admin — Dewata Nation Roleplay</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link href="https://fonts.googleapis.com/css2?family=Cinzel:wght@400;600;800&family=Crimson+Pro:ital,wght@0,300;0,400;0,600;1,300;1,400&display=swap" rel="stylesheet">
  <style>
    :root {
      --gold:       #c9a84c;
      --gold-light: #e8c97a;
      --gold-dim:   #7a6130;
      --cream:      #f5eed9;
      --dark:       #0d0b07;
      --dark-mid:   #1a160e;
      --dark-panel: #15120a;
      --border:     rgba(201,168,76,0.25);
      --border-glow:rgba(201,168,76,0.6);
      --red-warn:   #c0392b;
    }

    *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

    html { scroll-behavior: smooth; }

    body {
      background-color: var(--dark);
      color: var(--cream);
      font-family: 'Crimson Pro', Georgia, serif;
      font-size: 17px;
      min-height: 100vh;
      overflow-x: hidden;
    }

    /* ── BACKGROUND ── */
    body::before {
      content: '';
      position: fixed; inset: 0;
      background:
        radial-gradient(ellipse 80% 50% at 50% 0%, rgba(201,168,76,0.08) 0%, transparent 70%),
        radial-gradient(ellipse 60% 60% at 80% 80%, rgba(120,40,10,0.12) 0%, transparent 60%),
        repeating-linear-gradient(
          0deg,
          transparent,
          transparent 40px,
          rgba(201,168,76,0.015) 40px,
          rgba(201,168,76,0.015) 41px
        ),
        repeating-linear-gradient(
          90deg,
          transparent,
          transparent 40px,
          rgba(201,168,76,0.015) 40px,
          rgba(201,168,76,0.015) 41px
        );
      pointer-events: none;
      z-index: 0;
    }

    /* ── HEADER ── */
    header {
      position: relative;
      text-align: center;
      padding: 60px 20px 48px;
      z-index: 1;
    }

    .header-ornament {
      font-size: 11px;
      letter-spacing: 0.4em;
      color: var(--gold-dim);
      text-transform: uppercase;
      margin-bottom: 16px;
    }

    .header-ornament span {
      display: inline-block;
      width: 40px;
      height: 1px;
      background: var(--gold-dim);
      vertical-align: middle;
      margin: 0 10px;
    }

    h1 {
      font-family: 'Cinzel', serif;
      font-size: clamp(26px, 5vw, 48px);
      font-weight: 800;
      color: var(--gold-light);
      letter-spacing: 0.05em;
      line-height: 1.15;
      text-shadow: 0 0 40px rgba(201,168,76,0.4);
    }

    h1 em {
      display: block;
      font-style: normal;
      font-size: 0.55em;
      font-weight: 400;
      color: var(--gold-dim);
      letter-spacing: 0.25em;
      margin-top: 6px;
    }

    .header-line {
      width: 220px;
      height: 1px;
      background: linear-gradient(90deg, transparent, var(--gold), transparent);
      margin: 24px auto 0;
    }

    /* ── CONTAINER ── */
    .container {
      position: relative;
      z-index: 1;
      max-width: 800px;
      margin: 0 auto;
      padding: 0 20px 80px;
    }

    /* ── INTRO BOX ── */
    .intro-box {
      background: var(--dark-panel);
      border: 1px solid var(--border);
      border-radius: 2px;
      padding: 28px 32px;
      margin-bottom: 40px;
      position: relative;
    }
    .intro-box::before {
      content: '';
      position: absolute;
      top: -1px; left: 40px; right: 40px;
      height: 2px;
      background: linear-gradient(90deg, transparent, var(--gold), transparent);
    }
    .intro-box p {
      color: rgba(245,238,217,0.75);
      line-height: 1.75;
      font-size: 15px;
    }
    .intro-box strong { color: var(--gold-light); }

    /* ── SECTION ── */
    .section {
      margin-bottom: 36px;
      background: var(--dark-panel);
      border: 1px solid var(--border);
      border-radius: 2px;
      overflow: hidden;
      opacity: 0;
      transform: translateY(20px);
      animation: fadeUp 0.5s ease forwards;
    }
    .section:nth-child(1) { animation-delay: 0.1s; }
    .section:nth-child(2) { animation-delay: 0.2s; }
    .section:nth-child(3) { animation-delay: 0.3s; }
    .section:nth-child(4) { animation-delay: 0.4s; }
    .section:nth-child(5) { animation-delay: 0.5s; }

    @keyframes fadeUp {
      to { opacity: 1; transform: translateY(0); }
    }

    .section-header {
      display: flex;
      align-items: center;
      gap: 14px;
      padding: 18px 24px;
      background: rgba(201,168,76,0.06);
      border-bottom: 1px solid var(--border);
    }

    .section-icon {
      width: 32px; height: 32px;
      border: 1px solid var(--gold-dim);
      border-radius: 50%;
      display: flex; align-items: center; justify-content: center;
      font-size: 14px;
      flex-shrink: 0;
      color: var(--gold);
    }

    .section-title {
      font-family: 'Cinzel', serif;
      font-size: 13px;
      letter-spacing: 0.15em;
      color: var(--gold);
      text-transform: uppercase;
    }

    .section-subtitle {
      font-size: 12px;
      color: var(--gold-dim);
      margin-top: 2px;
    }

    .section-body {
      padding: 24px;
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 20px;
    }

    .field-full { grid-column: 1 / -1; }

    /* ── FIELD ── */
    .field label {
      display: block;
      font-size: 11px;
      letter-spacing: 0.12em;
      text-transform: uppercase;
      color: var(--gold-dim);
      margin-bottom: 8px;
    }
    .field label .req { color: var(--gold); margin-left: 2px; }

    .field input,
    .field select,
    .field textarea {
      width: 100%;
      background: rgba(255,255,255,0.03);
      border: 1px solid rgba(201,168,76,0.2);
      border-radius: 2px;
      padding: 10px 14px;
      color: var(--cream);
      font-family: 'Crimson Pro', serif;
      font-size: 15px;
      transition: border-color 0.2s, background 0.2s, box-shadow 0.2s;
      outline: none;
      -webkit-appearance: none;
    }

    .field select {
      cursor: pointer;
      background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='10' height='6' fill='%23c9a84c'%3E%3Cpath d='M0 0l5 6 5-6z'/%3E%3C/svg%3E");
      background-repeat: no-repeat;
      background-position: right 14px center;
      padding-right: 36px;
    }

    .field textarea {
      resize: vertical;
      min-height: 100px;
      line-height: 1.6;
    }

    .field input:focus,
    .field select:focus,
    .field textarea:focus {
      border-color: var(--gold);
      background: rgba(201,168,76,0.05);
      box-shadow: 0 0 0 3px rgba(201,168,76,0.08);
    }

    .field input::placeholder,
    .field textarea::placeholder {
      color: rgba(245,238,217,0.2);
      font-style: italic;
    }

    /* ── SUBMIT ── */
    .submit-area {
      text-align: center;
      margin-top: 40px;
    }

    .btn-submit {
      font-family: 'Cinzel', serif;
      font-size: 13px;
      letter-spacing: 0.2em;
      text-transform: uppercase;
      color: var(--dark);
      background: linear-gradient(135deg, var(--gold-light), var(--gold));
      border: none;
      padding: 16px 56px;
      cursor: pointer;
      border-radius: 2px;
      position: relative;
      overflow: hidden;
      transition: transform 0.2s, box-shadow 0.2s;
      box-shadow: 0 4px 24px rgba(201,168,76,0.25);
    }

    .btn-submit:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 32px rgba(201,168,76,0.4);
    }

    .btn-submit:active { transform: translateY(0); }

    .btn-submit:disabled {
      opacity: 0.6;
      cursor: not-allowed;
      transform: none;
    }

    .btn-submit .btn-text { position: relative; z-index: 1; }

    /* ── LOADER ── */
    .loader {
      display: none;
      width: 16px; height: 16px;
      border: 2px solid rgba(13,11,7,0.3);
      border-top-color: var(--dark);
      border-radius: 50%;
      animation: spin 0.7s linear infinite;
      margin: 0 auto;
    }
    @keyframes spin { to { transform: rotate(360deg); } }

    /* ── TOAST ── */
    .toast {
      position: fixed;
      bottom: 32px; left: 50%;
      transform: translateX(-50%) translateY(20px);
      background: var(--dark-panel);
      border: 1px solid var(--border);
      padding: 16px 28px;
      border-radius: 2px;
      font-size: 14px;
      letter-spacing: 0.03em;
      opacity: 0;
      transition: opacity 0.3s, transform 0.3s;
      z-index: 999;
      text-align: center;
      max-width: 380px;
    }
    .toast.show {
      opacity: 1;
      transform: translateX(-50%) translateY(0);
    }
    .toast.success { border-color: var(--gold); color: var(--gold-light); }
    .toast.error   { border-color: var(--red-warn); color: #e74c3c; }

    /* ── FOOTER ── */
    footer {
      text-align: center;
      padding: 24px;
      font-size: 12px;
      letter-spacing: 0.08em;
      color: var(--gold-dim);
      border-top: 1px solid var(--border);
      position: relative; z-index: 1;
    }

    /* ── RESPONSIVE ── */
    @media (max-width: 580px) {
      .section-body { grid-template-columns: 1fr; }
      .field-full   { grid-column: 1; }
      h1 { font-size: 24px; }
    }
  </style>
</head>
<body>

<header>
  <p class="header-ornament"><span></span>Rekrutmen Resmi<span></span></p>
  <h1>Dewata Nation Roleplay<em>Pendaftaran Tim Administrator</em></h1>
  <div class="header-line"></div>
</header>

<div class="container">

  <div class="intro-box">
    <p>
      Selamat datang di portal pendaftaran <strong>Administrator Dewata Nation Roleplay</strong>.
      Isi seluruh formulir dengan jujur dan lengkap. Setelah dikirim, lamaran kamu akan
      diteruskan langsung ke tim seleksi via WhatsApp. Pastikan data yang kamu berikan
      <strong>valid dan dapat diverifikasi</strong>.
    </p>
  </div>

  <form id="adminForm" novalidate>

    <!-- ── SECTION 1: OOC ── -->
    <div class="section">
      <div class="section-header">
        <div class="section-icon">👤</div>
        <div>
          <div class="section-title">Informasi OOC</div>
          <div class="section-subtitle">Out of Character — Data diri asli kamu</div>
        </div>
      </div>
      <div class="section-body">
        <div class="field">
          <label>Nama Asli <span class="req">*</span></label>
          <input type="text" name="ooc_nama_asli" placeholder="Nama lengkap kamu" required />
        </div>
        <div class="field">
          <label>Umur <span class="req">*</span></label>
          <input type="number" name="ooc_umur" placeholder="Usia kamu (tahun)" min="15" max="60" required />
        </div>
        <div class="field">
          <label>Wilayah / Kota <span class="req">*</span></label>
          <input type="text" name="ooc_wilayah" placeholder="Kota atau provinsi tempat tinggal" required />
        </div>
        <div class="field">
          <label>Nomor HP / WhatsApp <span class="req">*</span></label>
          <input type="tel" name="ooc_no_hp" placeholder="08xxxxxxxxxx" required />
        </div>
        <div class="field field-full">
          <label>Username Discord <span class="req">*</span></label>
          <input type="text" name="ooc_discord" placeholder="username#0000 atau username baru" required />
        </div>
        <div class="field field-full">
          <label>Sudah bergabung di server berapa lama? <span class="req">*</span></label>
          <select name="ooc_sudah_berapa" required>
            <option value="" disabled selected>— Pilih durasi —</option>
            <option value="Kurang dari 1 bulan">Kurang dari 1 bulan</option>
            <option value="1–3 bulan">1–3 bulan</option>
            <option value="3–6 bulan">3–6 bulan</option>
            <option value="6–12 bulan">6–12 bulan</option>
            <option value="Lebih dari 1 tahun">Lebih dari 1 tahun</option>
          </select>
        </div>
      </div>
    </div>

    <!-- ── SECTION 2: IC ── -->
    <div class="section">
      <div class="section-header">
        <div class="section-icon">🎭</div>
        <div>
          <div class="section-title">Informasi IC</div>
          <div class="section-subtitle">In Character — Data karakter di dalam game</div>
        </div>
      </div>
      <div class="section-body">
        <div class="field">
          <label>Nama Karakter <span class="req">*</span></label>
          <input type="text" name="ic_nama_karakter" placeholder="Nama karakter IC kamu" required />
        </div>
        <div class="field">
          <label>Umur Karakter <span class="req">*</span></label>
          <input type="number" name="ic_umur_karakter" placeholder="Umur karakter (IC)" min="17" max="80" required />
        </div>
        <div class="field">
          <label>Pekerjaan IC <span class="req">*</span></label>
          <input type="text" name="ic_pekerjaan_ic" placeholder="Pekerjaan karakter di server" required />
        </div>
        <div class="field">
          <label>Level Karakter <span class="req">*</span></label>
          <input type="number" name="ic_level_char" placeholder="Level karakter saat ini" min="1" required />
        </div>
        <div class="field field-full">
          <label>Pernah di-warn atau di-ban? <span class="req">*</span></label>
          <select name="ic_warned_banned" required>
            <option value="" disabled selected>— Pilih jawaban —</option>
            <option value="Tidak pernah">Tidak pernah</option>
            <option value="Pernah di-warn, sudah selesai">Pernah di-warn, sudah selesai</option>
            <option value="Pernah di-ban, sudah di-unban">Pernah di-ban, sudah di-unban</option>
            <option value="Sedang aktif warn/ban">Sedang aktif warn/ban</option>
          </select>
        </div>
      </div>
    </div>

    <!-- ── SECTION 3: PENGALAMAN ── -->
    <div class="section">
      <div class="section-header">
        <div class="section-icon">⚙️</div>
        <div>
          <div class="section-title">Pengalaman Admin</div>
          <div class="section-subtitle">Riwayat dan keahlian administrasi server</div>
        </div>
      </div>
      <div class="section-body">
        <div class="field field-full">
          <label>Apakah kamu pernah menjadi admin di server lain?</label>
          <select name="pengalaman_admin">
            <option value="Belum pernah">Belum pernah</option>
            <option value="Pernah, di server kecil">Pernah, di server kecil</option>
            <option value="Pernah, di server medium">Pernah, di server medium</option>
            <option value="Pernah, di server besar">Pernah, di server besar</option>
          </select>
        </div>
        <div class="field field-full">
          <label>Sebutkan nama server sebelumnya (jika ada)</label>
          <input type="text" name="server_sebelumnya" placeholder="Nama server (ketik 'Tidak ada' jika belum pernah)" />
        </div>
        <div class="field">
          <label>Sudah berapa lama bermain SAMP? <span class="req">*</span></label>
          <select name="lama_bermain" required>
            <option value="" disabled selected>— Pilih —</option>
            <option value="Kurang dari 6 bulan">Kurang dari 6 bulan</option>
            <option value="6 bulan – 1 tahun">6 bulan – 1 tahun</option>
            <option value="1–2 tahun">1–2 tahun</option>
            <option value="2–3 tahun">2–3 tahun</option>
            <option value="Lebih dari 3 tahun">Lebih dari 3 tahun</option>
          </select>
        </div>
        <div class="field">
          <label>Keahlian Khusus (opsional)</label>
          <input type="text" name="keahlian_khusus" placeholder="Mis: scripting, desain, moderasi" />
        </div>
      </div>
    </div>

    <!-- ── SECTION 4: MOTIVASI ── -->
    <div class="section">
      <div class="section-header">
        <div class="section-icon">✍️</div>
        <div>
          <div class="section-title">Motivasi & Komitmen</div>
          <div class="section-subtitle">Ceritakan niat dan rencana kontribusimu</div>
        </div>
      </div>
      <div class="section-body">
        <div class="field field-full">
          <label>Mengapa kamu ingin menjadi admin Dewata Nation RP? <span class="req">*</span></label>
          <textarea name="motivasi" placeholder="Tuliskan motivasi kamu dengan jelas dan jujur..." required></textarea>
        </div>
        <div class="field field-full">
          <label>Apa yang akan kamu kontribusikan untuk server ini? <span class="req">*</span></label>
          <textarea name="kontribusi_rencana" placeholder="Ide, rencana, atau komitmen nyata yang ingin kamu wujudkan..." required></textarea>
        </div>
        <div class="field">
          <label>Ketersediaan Waktu per Hari <span class="req">*</span></label>
          <select name="ketersediaan_waktu" required>
            <option value="" disabled selected>— Pilih —</option>
            <option value="1–2 jam/hari">1–2 jam/hari</option>
            <option value="2–4 jam/hari">2–4 jam/hari</option>
            <option value="4–6 jam/hari">4–6 jam/hari</option>
            <option value="Lebih dari 6 jam/hari">Lebih dari 6 jam/hari</option>
          </select>
        </div>
        <div class="field">
          <label>Skenario: Ada pemain lapor cheater, kamu tidak online. Apa yang kamu lakukan? <span class="req">*</span></label>
          <textarea name="skenario_handler" placeholder="Jawab dengan singkat dan jelas..." style="min-height:80px;" required></textarea>
        </div>
      </div>
    </div>

    <!-- ── SUBMIT ── -->
    <div class="submit-area">
      <button type="submit" class="btn-submit" id="submitBtn">
        <span class="btn-text">Kirim Lamaran</span>
        <div class="loader" id="loader"></div>
      </button>
      <p style="margin-top:16px; font-size:13px; color:rgba(245,238,217,0.35); letter-spacing:0.04em;">
        Data akan diteruskan ke tim rekrutmen via WhatsApp
      </p>
    </div>

  </form>
</div>

<footer>
  &copy; 2025 Dewata Nation Roleplay &mdash; All rights reserved
</footer>

<div class="toast" id="toast"></div>

<script>
  function showToast(msg, type) {
    const t = document.getElementById('toast');
    t.textContent = msg;
    t.className = 'toast ' + type + ' show';
    setTimeout(() => t.classList.remove('show'), 4500);
  }

  document.getElementById('adminForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    const form = e.target;
    const btn  = document.getElementById('submitBtn');
    const loader = document.getElementById('loader');

    // Basic validation
    const required = form.querySelectorAll('[required]');
    let valid = true;
    required.forEach(el => {
      if (!el.value.trim()) {
        el.style.borderColor = '#c0392b';
        valid = false;
      } else {
        el.style.borderColor = '';
      }
    });
    if (!valid) { showToast('Mohon lengkapi semua field yang wajib diisi.', 'error'); return; }

    const data = {};
    new FormData(form).forEach((v, k) => data[k] = v);

    btn.disabled = true;
    btn.querySelector('.btn-text').style.display = 'none';
    loader.style.display = 'block';

    try {
      const res = await fetch('/submit', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
      });
      const json = await res.json();
      if (json.success) {
        showToast('✓ Lamaran berhasil dikirim! Tim kami akan segera meninjau.', 'success');
        form.reset();
      } else {
        showToast('Gagal mengirim: ' + json.message, 'error');
      }
    } catch(err) {
      showToast('Terjadi kesalahan koneksi. Coba lagi.', 'error');
    } finally {
      btn.disabled = false;
      btn.querySelector('.btn-text').style.display = '';
      loader.style.display = 'none';
    }
  });
</script>
</body>
</html>`

// ─── HANDLERS ─────────────────────────────────────────────────────────────────

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "Template error", 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(Response{false, "Method not allowed"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(Response{false, "Gagal membaca data"})
		return
	}
	defer r.Body.Close()

	var form AdminForm
	if err := json.Unmarshal(body, &form); err != nil {
		json.NewEncoder(w).Encode(Response{false, "Data tidak valid"})
		return
	}

	message := buildWhatsAppMessage(form)

	if err := sendWhatsApp(message); err != nil {
		log.Printf("WhatsApp error: %v", err)
		json.NewEncoder(w).Encode(Response{false, "Gagal mengirim ke WhatsApp: " + err.Error()})
		return
	}

	json.NewEncoder(w).Encode(Response{true, "Lamaran berhasil dikirim"})
}

// ─── FORMAT PESAN ─────────────────────────────────────────────────────────────

func buildWhatsAppMessage(f AdminForm) string {
	now := time.Now().Format("02 Jan 2006, 15:04 WIB")

	var sb strings.Builder

	sb.WriteString("╔══════════════════════════╗\n")
	sb.WriteString("║  PENDAFTARAN ADMIN        ║\n")
	sb.WriteString("║  DEWATA NATION ROLEPLAY   ║\n")
	sb.WriteString("╚══════════════════════════╝\n")
	sb.WriteString(fmt.Sprintf("📅 Diterima: %s\n\n", now))

	sb.WriteString("━━━ 👤 DATA OOC (Out of Character) ━━━\n")
	sb.WriteString(fmt.Sprintf("• Nama Asli         : %s\n", f.OOCNamaAsli))
	sb.WriteString(fmt.Sprintf("• Umur              : %s tahun\n", f.OOCUmur))
	sb.WriteString(fmt.Sprintf("• Wilayah           : %s\n", f.OOCWilayah))
	sb.WriteString(fmt.Sprintf("• No. HP/WA         : %s\n", f.OOCNoHP))
	sb.WriteString(fmt.Sprintf("• Discord           : %s\n", f.OOCDiscord))
	sb.WriteString(fmt.Sprintf("• Lama di Server    : %s\n\n", f.OOCSudahBerapa))

	sb.WriteString("━━━ 🎭 DATA IC (In Character) ━━━\n")
	sb.WriteString(fmt.Sprintf("• Nama Karakter     : %s\n", f.ICNamaKarakter))
	sb.WriteString(fmt.Sprintf("• Umur Karakter     : %s tahun\n", f.ICUmurKarakter))
	sb.WriteString(fmt.Sprintf("• Pekerjaan IC      : %s\n", f.ICPekerjaanIC))
	sb.WriteString(fmt.Sprintf("• Level Karakter    : %s\n", f.ICLevelChar))
	sb.WriteString(fmt.Sprintf("• Riwayat Warn/Ban  : %s\n\n", f.ICWarnedBanned))

	sb.WriteString("━━━ ⚙️ PENGALAMAN ADMIN ━━━\n")
	sb.WriteString(fmt.Sprintf("• Pengalaman Admin  : %s\n", f.PengalamanAdmin))
	sb.WriteString(fmt.Sprintf("• Server Sebelumnya : %s\n", f.ServerSebelumnya))
	sb.WriteString(fmt.Sprintf("• Lama Main SAMP    : %s\n", f.LamaBermain))
	sb.WriteString(fmt.Sprintf("• Keahlian Khusus   : %s\n\n", f.KeahlianKhusus))

	sb.WriteString("━━━ ✍️ MOTIVASI & KOMITMEN ━━━\n")
	sb.WriteString(fmt.Sprintf("• Motivasi          :\n  %s\n\n", wordWrap(f.Motivasi, "  ")))
	sb.WriteString(fmt.Sprintf("• Rencana Kontribusi:\n  %s\n\n", wordWrap(f.KontribusiRencana, "  ")))
	sb.WriteString(fmt.Sprintf("• Waktu Tersedia    : %s\n", f.KetersediaanWaktu))
	sb.WriteString(fmt.Sprintf("• Skenario Admin    :\n  %s\n\n", wordWrap(f.SkenarioHandler, "  ")))

	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	sb.WriteString("_Pesan otomatis dari portal rekrutmen Dewata Nation RP_")

	return sb.String()
}

func wordWrap(text, prefix string) string {
	// Replace newlines for clean formatting
	replaced := strings.ReplaceAll(strings.TrimSpace(text), "\n", "\n"+prefix)
	return replaced
}

// ─── KIRIM WHATSAPP ───────────────────────────────────────────────────────────

func sendWhatsApp(message string) error {
	token := getFonnteToken()

	formData := url.Values{}
	formData.Set("target", WhatsAppTarget)
	formData.Set("message", message)

	req, err := http.NewRequest("POST", FonnteAPIURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("Fonnte response [%d]: %s", resp.StatusCode, string(respBody))

	if resp.StatusCode != 200 {
		return fmt.Errorf("fonnte returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// ─── MAIN ─────────────────────────────────────────────────────────────────────

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/submit", submitHandler)

	port := getPort()
	log.Printf("🚀 Dewata Nation RP Admin Portal running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}