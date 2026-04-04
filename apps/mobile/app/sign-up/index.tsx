import {SafeAreaView, StyleSheet, Dimensions} from "react-native";
import {useForm, SubmitHandler} from 'react-hook-form';
import ControlledTextInput from "@/components/FormControls/ControlledTextInput";
import {Button, Card, Text} from 'react-native-paper'
import Joi from "joi";
import {joiResolver} from "@hookform/resolvers/joi";
import {useFormSubmit} from "@/hooks/useFormSubmit";

interface SignUpData {
    email: string;
    password: string;
}

const schema = Joi.object({
    email: Joi.string()
        .email({tlds: {allow: false}})
        .required()
        .messages({
            "string.empty": "Email is required",
            "string.email": "Please enter a valid email address",
        }),
    password: Joi.string()
        .min(6)
        .required()
        .messages({
            "string.empty": "Password is required",
            "string.min": "Password must be at least 6 characters",
        }),
});

export default function SignUp() {
    const {control, handleSubmit, formState: {errors}} = useForm<SignUpData>(
        {
            resolver: joiResolver(schema),
        }
    );
    const {submitForm, error, loading} = useFormSubmit();
    const onSubmit: SubmitHandler<SignUpData> = async (data) => {
        await submitForm('/v1/user/register', data);
    }

    return (
        <SafeAreaView style={styles.container}>
            <Card style={styles.card}>
                <Card.Title title={'Sign-Up'}/>
                <Card.Content>
                    {error && <Text style={{ color: 'red' }}>{error}</Text>}
                    <ControlledTextInput<SignUpData> style={styles.formElement} control={control} label='Email'
                                                     name={'email'} error={errors.email?.message}/>
                    <ControlledTextInput<SignUpData> style={styles.formElement} control={control} label='Password'
                                                     name={'password'} secureTextEntry
                                                     error={errors.password?.message}/>
                </Card.Content>
                <Card.Actions>
                    <Button style={[styles.submitButton]} onPress={handleSubmit(onSubmit)}>Sign-in</Button>
                </Card.Actions>
            </Card>
        </SafeAreaView>
    )
}


const windowHeight = Dimensions.get('window').height;
const windowWidth = Dimensions.get('window').width;

const styles = StyleSheet.create({
    container: {
        flex: 1,
        marginTop: windowHeight * .30,
        alignItems: 'center',
    },
    submitButton: {
        alignSelf: 'flex-end'
    },
    formElement: {
        marginTop: 8,
    },
    card: {
        maxWidth: 700,
        width: '90%'
    }
});
