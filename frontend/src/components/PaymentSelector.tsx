import React from 'react';
import { useTranslation } from 'react-i18next';
import { motion } from 'framer-motion';
import { cn } from '@/lib/utils';

export interface PaymentGateway {
  id: 'momo' | 'zalopay' | 'payos';
  name: string;
  logo: string;
  description: string;
  color: string;
}

interface PaymentSelectorProps {
  selectedGateway?: string;
  onSelect: (gateway: string) => void;
  className?: string;
}

export const PaymentSelector: React.FC<PaymentSelectorProps> = ({
  selectedGateway,
  onSelect,
  className
}) => {
  const { t } = useTranslation();

  const gateways: PaymentGateway[] = [
    {
      id: 'momo',
      name: 'MoMo',
      logo: '💳',
      description: t('payment.methods.momo_description'),
      color: 'from-pink-500 to-rose-500'
    },
    {
      id: 'zalopay',
      name: 'ZaloPay',
      logo: '🔵',
      description: t('payment.methods.zalopay_description'),
      color: 'from-blue-500 to-cyan-500'
    },
    {
      id: 'payos',
      name: 'PayOS',
      logo: '💰',
      description: t('payment.methods.bank_description'),
      color: 'from-emerald-500 to-teal-500'
    }
  ];

  return (
    <div className={cn('space-y-4', className)}>
      <h3 className="text-lg font-semibold text-gray-900">
        {t('payment.select_method')}
      </h3>
      
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {gateways.map((gateway) => (
          <motion.button
            key={gateway.id}
            type="button"
            onClick={() => onSelect(gateway.id)}
            className={cn(
              'relative p-6 rounded-xl border-2 transition-all text-left',
              selectedGateway === gateway.id
                ? 'border-blue-500 bg-blue-50'
                : 'border-gray-200 hover:border-gray-300 bg-white'
            )}
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
          >
            {/* Selection indicator */}
            {selectedGateway === gateway.id && (
              <motion.div
                layoutId="selected-payment"
                className="absolute inset-0 bg-gradient-to-br from-blue-100 to-indigo-100 rounded-xl"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.2 }}
              />
            )}

            <div className="relative z-10">
              {/* Logo */}
              <div className={cn(
                'w-16 h-16 rounded-lg bg-gradient-to-br flex items-center justify-center mb-3',
                gateway.color
              )}>
                <span className="text-3xl">{gateway.logo}</span>
              </div>

              {/* Name */}
              <h4 className="text-lg font-semibold text-gray-900 mb-1">
                {gateway.name}
              </h4>

              {/* Description */}
              <p className="text-sm text-gray-600">
                {gateway.description}
              </p>

              {/* Selected checkmark */}
              {selectedGateway === gateway.id && (
                <motion.div
                  initial={{ scale: 0 }}
                  animate={{ scale: 1 }}
                  className="absolute top-4 right-4 w-6 h-6 bg-blue-500 rounded-full flex items-center justify-center"
                >
                  <svg className="w-4 h-4 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                </motion.div>
              )}
            </div>
          </motion.button>
        ))}
      </div>
    </div>
  );
};

interface PaymentButtonProps {
  isLoading: boolean;
  disabled: boolean;
  onClick: () => void;
  amount: number;
}

export const PaymentButton: React.FC<PaymentButtonProps> = ({
  isLoading,
  disabled,
  onClick,
  amount
}) => {
  const { t } = useTranslation();

  return (
    <motion.button
      type="button"
      onClick={onClick}
      disabled={disabled || isLoading}
      className={cn(
        'w-full py-4 px-6 rounded-xl font-semibold text-white transition-all',
        'bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700',
        'disabled:from-gray-400 disabled:to-gray-500 disabled:cursor-not-allowed',
        'shadow-lg hover:shadow-xl'
      )}
      whileHover={{ scale: disabled ? 1 : 1.02 }}
      whileTap={{ scale: disabled ? 1 : 0.98 }}
    >
      {isLoading ? (
        <span className="flex items-center justify-center gap-2">
          <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
          {t('payment.processing')}
        </span>
      ) : (
        <span className="flex items-center justify-center gap-2">
          {t('payment.proceed_to_payment')} • {amount.toLocaleString('vi-VN')} ₫
        </span>
      )}
    </motion.button>
  );
};

interface PaymentStatusProps {
  status: 'pending' | 'processing' | 'completed' | 'failed' | 'cancelled';
  message?: string;
}

export const PaymentStatus: React.FC<PaymentStatusProps> = ({ status, message }) => {
  const { t } = useTranslation();

  const statusConfig = {
    pending: {
      icon: '⏳',
      text: t('payment.status.pending'),
      color: 'bg-yellow-100 text-yellow-800 border-yellow-300'
    },
    processing: {
      icon: '🔄',
      text: t('payment.status.processing'),
      color: 'bg-blue-100 text-blue-800 border-blue-300'
    },
    completed: {
      icon: '✅',
      text: t('payment.status.completed'),
      color: 'bg-green-100 text-green-800 border-green-300'
    },
    failed: {
      icon: '❌',
      text: t('payment.status.failed'),
      color: 'bg-red-100 text-red-800 border-red-300'
    },
    cancelled: {
      icon: '🚫',
      text: t('payment.status.cancelled'),
      color: 'bg-gray-100 text-gray-800 border-gray-300'
    }
  };

  const config = statusConfig[status];

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className={cn(
        'p-4 rounded-lg border-2 flex items-center gap-3',
        config.color
      )}
    >
      <span className="text-2xl">{config.icon}</span>
      <div className="flex-1">
        <p className="font-semibold">{config.text}</p>
        {message && <p className="text-sm mt-1">{message}</p>}
      </div>
    </motion.div>
  );
};
