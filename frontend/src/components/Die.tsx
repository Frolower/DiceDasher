export function Die({ className = '', value = '20' }: { className?: string; value?: string }) {
  return (
    <svg className={className} viewBox="0 0 320 350" fill="none" aria-hidden="true">
      <polygon points="160,12 296,91 291,254 160,335 25,254 24,92" fill="currentColor" />
      <polygon points="160,12 96,114 24,92" fill="#fff" fillOpacity=".36" />
      <polygon points="160,12 296,91 224,114" fill="#fff" fillOpacity=".18" />
      <polygon points="24,92 96,114 25,254" fill="#fff" fillOpacity=".2" />
      <polygon points="296,91 224,114 291,254" fill="#001a61" fillOpacity=".12" />
      <polygon points="96,114 224,114 160,234" fill="#fff" fillOpacity=".16" />
      <polygon points="25,254 160,234 160,335" fill="#001a61" fillOpacity=".16" />
      <polygon points="291,254 160,234 160,335" fill="#001a61" fillOpacity=".28" />
      <path
        d="M160 12 96 114 24 92M160 12 224 114 296 91M96 114H224L160 234 96 114 25 254 160 234 291 254 224 114M160 234V335"
        stroke="white"
        strokeOpacity=".53"
        strokeWidth="1.5"
      />
      <text
        x="160"
        y="170"
        textAnchor="middle"
        dominantBaseline="middle"
        fill="white"
        fontSize="53"
        fontFamily="Georgia, serif"
      >
        {value}
      </text>
    </svg>
  )
}
