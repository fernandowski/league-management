import { joiResolver } from '@hookform/resolvers/joi';
import Joi from 'joi';
import { SubmitHandler, useForm } from 'react-hook-form';
import { StyleSheet, View } from 'react-native';

import ControlledTextInput from '@/components/FormControls/ControlledTextInput';
import { AppButton } from '@/components/ui/AppButton';
import { AppCard } from '@/components/ui/AppCard';
import { AppScreen } from '@/components/ui/AppScreen';
import { AppText } from '@/components/ui/AppText';
import { useFormSubmit } from '@/hooks/useFormSubmit';
import { themeTokens, useAppTheme } from '@/theme/theme';

interface SignUpData {
  email: string;
  password: string;
}

const schema = Joi.object({
  email: Joi.string()
    .email({ tlds: { allow: false } })
    .required()
    .messages({
      'string.empty': 'Email is required',
      'string.email': 'Please enter a valid email address',
    }),
  password: Joi.string()
    .min(6)
    .required()
    .messages({
      'string.empty': 'Password is required',
      'string.min': 'Password must be at least 6 characters',
    }),
});

export default function SignUp() {
  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<SignUpData>({
    resolver: joiResolver(schema),
    defaultValues: {
      email: '',
      password: '',
    },
  });
  const { submitForm, error } = useFormSubmit();
  const theme = useAppTheme();

  const onSubmit: SubmitHandler<SignUpData> = async (data) => {
    await submitForm('/v1/user/register', data);
  };

  return (
    <AppScreen style={styles.screen}>
      <View style={styles.accentBand} />
      <AppCard style={styles.card}>
        <AppCard.Content style={styles.content}>
          <AppText variant="headlineMedium" style={styles.title}>
            Create Your Account
          </AppText>
          <AppText variant="bodyLarge" style={[styles.subtitle, { color: theme.colors.onSurfaceVariant }]}>
            Set up your workspace and start running leagues with a consistent admin flow.
          </AppText>
          {error && <AppText style={{ color: theme.colors.error }}>{error}</AppText>}
          <ControlledTextInput<SignUpData>
            style={styles.formElement}
            control={control}
            label="Email"
            name="email"
            error={errors.email?.message}
            autoCapitalize="none"
          />
          <ControlledTextInput<SignUpData>
            style={styles.formElement}
            control={control}
            label="Password"
            name="password"
            secureTextEntry
            error={errors.password?.message}
          />
          <AppButton mode="contained" style={styles.submitButton} onPress={handleSubmit(onSubmit)}>
            Sign Up
          </AppButton>
        </AppCard.Content>
      </AppCard>
    </AppScreen>
  );
}

const styles = StyleSheet.create({
  screen: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 20,
  },
  accentBand: {
    position: 'absolute',
    bottom: 80,
    left: -40,
    width: 260,
    height: 180,
    borderRadius: 36,
    transform: [{ rotate: '-12deg' }],
    backgroundColor: themeTokens.colors.primarySoft,
  },
  card: {
    width: '100%',
    maxWidth: 520,
    borderWidth: 1,
    borderRadius: 28,
  },
  content: {
    gap: 10,
    paddingVertical: 12,
  },
  title: {
    marginBottom: 6,
  },
  subtitle: {
    marginBottom: 10,
  },
  submitButton: {
    marginTop: 8,
    alignSelf: 'flex-start',
  },
  formElement: {
    marginTop: 4,
  },
});
