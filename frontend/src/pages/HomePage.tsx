import { Link } from 'react-router-dom'
import { Search, TrendingUp, Shield, Zap } from 'lucide-react'

export default function HomePage() {
  return (
    <div>
      {/* Hero Section */}
      <section className="bg-gradient-to-br from-primary-600 to-primary-800 text-white py-20">
        <div className="container-custom">
          <div className="max-w-3xl">
            <h1 className="text-5xl font-bold mb-6">
              Book Your Bus Tickets with Real-Time Seat Selection
            </h1>
            <p className="text-xl mb-8 text-primary-100">
              Fast, secure, and convenient bus booking platform with live seat availability
            </p>
            <Link to="/search" className="btn bg-white text-primary-600 hover:bg-gray-100 px-8 py-4 text-lg">
              <Search className="inline-block mr-2 h-5 w-5" />
              Search Trips
            </Link>
          </div>
        </div>
      </section>

      {/* Features */}
      <section className="py-16">
        <div className="container-custom">
          <h2 className="text-3xl font-bold text-center mb-12">Why Choose Us?</h2>
          <div className="grid md:grid-cols-3 gap-8">
            <div className="card p-6 text-center">
              <Zap className="h-12 w-12 text-primary-600 mx-auto mb-4" />
              <h3 className="text-xl font-semibold mb-2">Real-Time Updates</h3>
              <p className="text-gray-600">
                See live seat availability and get instant booking confirmations
              </p>
            </div>
            <div className="card p-6 text-center">
              <Shield className="h-12 w-12 text-primary-600 mx-auto mb-4" />
              <h3 className="text-xl font-semibold mb-2">Secure Payments</h3>
              <p className="text-gray-600">
                Multiple payment options with industry-standard security
              </p>
            </div>
            <div className="card p-6 text-center">
              <TrendingUp className="h-12 w-12 text-primary-600 mx-auto mb-4" />
              <h3 className="text-xl font-semibold mb-2">Best Prices</h3>
              <p className="text-gray-600">
                Competitive pricing with regular offers and discounts
              </p>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}
