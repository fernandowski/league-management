import {View, StyleSheet} from "react-native";
import TeamList from "@/components/Teams/TeamList";
import {useOrganizationStore} from "@/stores/organizationStore";
import {useCallback, useState} from "react";
import {apiRequest} from "@/api/api";
import {useFocusEffect} from "expo-router";
import ControlledTextInput from "@/components/FormControls/ControlledTextInput";
import Joi from "joi";
import {useForm} from "react-hook-form";
import {joiResolver} from "@hookform/resolvers/joi";
import ViewContent from "@/components/Layout/ViewContent";
import StyledModal from "@/components/StyledModal";
import { AppButton } from "@/components/ui/AppButton";
import { AppText } from "@/components/ui/AppText";

export interface Team {
    id: string
    name: string
    organizationId: string
}

export interface CreateTeamData {
    name: string
}

const schema = Joi.object({
    name: Joi.string()
        .required()
        .messages({
            "string.empty": "Team name is required",
        })
});

export default function TeamOverview() {
    const {organization} = useOrganizationStore();
    const [teams, setTeams] = useState<Team[]>([]);
    const [showModal, setShowModal] = useState(false);
    const {control, handleSubmit, formState: {errors}} = useForm<CreateTeamData>(
        {
            resolver: joiResolver(schema),
        }
    );
    const handleOpenModal = () => {
        setShowModal(!showModal);
    };

    const handleSave = async (data: CreateTeamData) => {
        try {
            await apiRequest('/v1/teams', {
                method: 'POST',
                body: {
                    name: data.name,
                    organization_id: organization
                }
            })
            fetchTeams();
            setShowModal(!showModal);
        } catch {

        }
    }

    const fetchTeams = useCallback(async () => {
        try {
            const response = await apiRequest(`/v1/teams/?organization_id=${organization}`, {method: 'GET'});
            setTeams(response.map((team: Record<string, string>) => ({
                id: team.id,
                name: team.name,
                organizationName: team.organization_name
            })));
        } catch {
            setTeams([]);
        }
    }, [organization])

    useFocusEffect(
        useCallback(() => {
            fetchTeams();
            return () => {
                setShowModal(false)
            }
        }, [fetchTeams])
    )

    return (
        <ViewContent>
            <AppButton style={[styles.addTeamButton]} mode={"elevated"} onPress={handleOpenModal}> + Add Team </AppButton>
            <TeamList data={teams}/>
            <StyledModal isOpen={showModal} onDismiss={handleOpenModal} contentContainerStyle={styles.modal}>
                <View style={[styles.formContainer]}>
                    <AppText variant="titleMedium">Team Name</AppText>
                    <ControlledTextInput label='Name' name={'name'} control={control} error={errors.name?.message}/>
                    <AppButton style={{alignSelf: "flex-end"}} onPress={handleSubmit(handleSave)}>Save</AppButton>
                </View>
            </StyledModal>
        </ViewContent>
    )
}


const styles = StyleSheet.create({
    addTeamButton: {
        alignSelf: "flex-end"
    },
    formContainer: {
        flex: 0.7,
        padding: 16
    },
    modal: {
        width: "80%",
        maxHeight: 400,
        maxWidth: 400,
        alignSelf: "center",
    }
})
