import {View, StyleSheet} from "react-native";
import {Button, Text, Modal} from "react-native-paper";
import {useState} from "react";
import ControlledTextInput from "@/components/FormControls/ControlledTextInput";
import Joi from "joi";
import {SubmitHandler, useForm} from "react-hook-form";
import {joiResolver} from "@hookform/resolvers/joi";
import {apiRequest} from "@/api/api";
import {useOrganizationStore} from "@/stores/organizationStore";


const schema = Joi.object({
    name: Joi.string()
        .required()
        .messages({
            "string.empty": "Organization name is required",
        })
});

interface CreateOrganizationData {
    name: string
}


export default function AddOrganization() {
    const [showModal, setShowModal] = useState<boolean>(false);
    const {fetchOrganizations} = useOrganizationStore();
    const {control, handleSubmit, formState: {errors}} = useForm<CreateOrganizationData>(
        {
            resolver: joiResolver(schema),
        }
    );

    const handleShowModal = () => {
        setShowModal(!showModal);
    }

    const handleOnSaveOrganization: SubmitHandler<CreateOrganizationData> = async (data: CreateOrganizationData): Promise<void> => {
        await apiRequest('/v1/organizations', {
            method: 'POST',
            body: data
        });
        fetchOrganizations();

        setShowModal(!showModal);
    }

    return (
        <View style={[styles.container]}>
            <Text>Create an organization to start managing your leagues</Text>
            <Button buttonColor={"#f194ff"} mode={"contained-tonal"} onPress={handleShowModal} dark>Add
                Organization</Button>
            <Modal visible={showModal} dismissable={false} contentContainerStyle={[styles.modal]}>
                <View style={[styles.formContainer]}>
                    <Text>Organization Name</Text>
                    <ControlledTextInput label='Name' name={'name'} control={control} error={errors.name?.message}/>
                    <Button style={{alignSelf: "flex-end"}} onPress={handleSubmit(handleOnSaveOrganization)}>Save</Button>
                </View>
            </Modal>
        </View>
    )
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        alignItems: "center",
        justifyContent: "center",
        gap: 8,
    },
    modal: {
        flex: 0.8,
        padding: 16,
        backgroundColor: "white",
        width: "80%",
        maxHeight: 400,
        maxWidth: 400,
        zIndex: 200,
        alignSelf: "center"
    },
    formContainer: {
        flex: 0.7,
        padding: 16
    }
})
