import {SafeAreaView, StyleSheet, Dimensions} from "react-native";
import {useForm, SubmitHandler} from 'react-hook-form';
import ControlledTextInput from "@/components/FormControls/ControlledTextInput";
import {Button, Card, Text} from 'react-native-paper'
import Joi from "joi";
import {joiResolver} from "@hookform/resolvers/joi";
import {useFormSubmit} from "@/hooks/useFormSubmit";
import {storeJWT} from "@/util/jwt-manager";
import {router} from "expo-router";

interface LoginData {
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
        }),
});

export default function Login() {
    const {control, handleSubmit, formState: {errors}} = useForm<LoginData>(
        {
            resolver: joiResolver(schema),
        }
    );
    const {submitForm, error, loading} = useFormSubmit();
    const onSubmit: SubmitHandler<LoginData> = async (data) => {
        const response = await submitForm('/v1/user/login', data);
        await storeJWT(response.jwt)
        router.push("/dashboard")
    }

    return (
        <SafeAreaView style={styles.container}>
            <Card style={styles.card}>
                <Card.Title title={'Login'}/>
                <Card.Content>
                    {error && <Text style={{ color: 'red' }}>{error}</Text>}
                    <ControlledTextInput<LoginData> style={styles.formElement} control={control} label='Email'
                                                     name={'email'} error={errors.email?.message}/>
                    <ControlledTextInput<LoginData> style={styles.formElement} control={control} label='Password'
                                                     name={'password'} secureTextEntry
                                                     error={errors.password?.message}/>
                </Card.Content>
                <Card.Actions>
                    <Button style={[styles.submitButton]} onPress={handleSubmit(onSubmit)}>Login</Button>
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
