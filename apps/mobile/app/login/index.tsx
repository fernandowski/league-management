import { joiResolver } from '@hookform/resolvers/joi';
import { router } from 'expo-router';
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
import { storeJWT } from '@/util/jwt-manager';

interface LoginData {
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
    }),
});

export default function Login() {
  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginData>({
    resolver: joiResolver(schema),
    defaultValues: {
      email: '',
      password: '',
    },
  });
  const { submitForm, error } = useFormSubmit();
  const theme = useAppTheme();

  const onSubmit: SubmitHandler<LoginData> = async (data) => {
    const response = await submitForm('/v1/user/login', data);
    await storeJWT(response.jwt);
    router.push('/dashboard');
  };

  return (
    <AppScreen style={styles.screen}>
      <View style={styles.accentOrb} />
      <AppCard style={styles.card}>
        <AppCard.Content style={styles.content}>
          <AppText variant="headlineMedium" style={styles.title}>
            League OS
          </AppText>
          <AppText variant="bodyLarge" style={[styles.subtitle, { color: theme.colors.onSurfaceVariant }]}>
            Sign in to manage organizations, leagues, and seasons from one place.
          </AppText>
          {error && <AppText style={{ color: theme.colors.error }}>{error}</AppText>}
          <ControlledTextInput<LoginData>
            style={styles.formElement}
            control={control}
            label="Email"
            name="email"
            error={errors.email?.message}
            autoCapitalize="none"
          />
          <ControlledTextInput<LoginData>
            style={styles.formElement}
            control={control}
            label="Password"
            name="password"
            secureTextEntry
            error={errors.password?.message}
          />
          <AppButton variant="submit" style={styles.submitButton} onPress={handleSubmit(onSubmit)}>
            Login
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
  accentOrb: {
    position: 'absolute',
    top: 72,
    right: 28,
    width: 180,
    height: 180,
    borderRadius: 90,
    backgroundColor: themeTokens.colors.secondarySoft,
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
