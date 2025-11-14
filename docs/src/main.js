import './style.css'
import Prism from 'prismjs'
import 'prismjs/components/prism-bash'
import 'prismjs/components/prism-yaml'
import 'prismjs/components/prism-sql'
import 'prismjs/themes/prism.css'

// Prism.jsの自動ハイライトを初期化
document.addEventListener('DOMContentLoaded', () => {
  Prism.highlightAll()
  initSlideNavigation()
})

// スライドナビゲーション機能の初期化
function initSlideNavigation() {
  const slides = document.querySelectorAll('.slide')
  const prevBtn = document.getElementById('prevSlideBtn')
  const nextBtn = document.getElementById('nextSlideBtn')
  const counter = document.getElementById('slideCounter')

  if (!slides.length || !prevBtn || !nextBtn || !counter) {
    return // スライド要素が存在しない場合は何もしない
  }

  let currentSlide = 0
  const totalSlides = slides.length

  function showSlide(index) {
    // インデックスの範囲チェック
    if (index < 0 || index >= totalSlides) {
      return
    }

    // 全スライドを非表示にし、指定されたスライドのみ表示
    slides.forEach((slide, i) => {
      if (i === index) {
        slide.classList.remove('hidden')
      } else {
        slide.classList.add('hidden')
      }
    })

    // カウンター更新
    counter.textContent = `${index + 1} / ${totalSlides}`

    // ボタンの有効/無効状態を更新
    prevBtn.disabled = index === 0
    nextBtn.disabled = index === totalSlides - 1

    // 現在のスライドインデックスを更新
    currentSlide = index
  }

  // 前へボタンのクリックイベント
  prevBtn.addEventListener('click', () => {
    if (currentSlide > 0) {
      showSlide(currentSlide - 1)
    }
  })

  // 次へボタンのクリックイベント
  nextBtn.addEventListener('click', () => {
    if (currentSlide < totalSlides - 1) {
      showSlide(currentSlide + 1)
    }
  })

  // キーボードナビゲーション（矢印キー）
  document.addEventListener('keydown', (e) => {
    // スライドセクションが表示されているかチェック
    const lectureSection = document.getElementById('lecture')
    if (!lectureSection) return

    // 座学セクションが画面内にあるかチェック（簡易的な判定）
    const rect = lectureSection.getBoundingClientRect()
    const isVisible = rect.top < window.innerHeight && rect.bottom > 0

    if (!isVisible) return

    if (e.key === 'ArrowLeft' && currentSlide > 0) {
      e.preventDefault()
      showSlide(currentSlide - 1)
    } else if (e.key === 'ArrowRight' && currentSlide < totalSlides - 1) {
      e.preventDefault()
      showSlide(currentSlide + 1)
    }
  })

  // 初期表示
  showSlide(0)
}
