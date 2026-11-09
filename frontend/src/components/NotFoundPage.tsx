import { useEffect, useRef } from 'react'
import { Link } from 'react-router'
import { t } from '../lib/i18n'
import { Die } from './Die'

export function NotFoundPage({ homePath }: { homePath: '/' | '/profile' }) {
  const heading = useRef<HTMLHeadingElement>(null)
  const title = t('text.pageNotFound')

  useEffect(() => {
    document.title = `404 · ${title} · DiceDasher`
    heading.current?.focus()
    return () => {
      document.title = 'DiceDasher'
    }
  }, [title])

  return (
    <main id="main" className="not-found-page shell" aria-labelledby="not-found-title">
      <div className="not-found-art" aria-hidden="true">
        <Die className="not-found-die" value="1" />
      </div>
      <div className="not-found-content">
        <p className="not-found-code">404</p>
        <h1 id="not-found-title" ref={heading} tabIndex={-1}>
          {title}
        </h1>
        <p className="not-found-description">{t('text.thisPageDoesnTExist')}</p>
        <Link className="button primary large" to={homePath}>
          <span aria-hidden="true">←</span>
          {t('text.goHome')}
        </Link>
      </div>
    </main>
  )
}
