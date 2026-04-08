import {useLeagueData} from "@/hooks/useLeagueData";
import {useOrganizationStore} from "@/stores/organizationStore";
import {StyleSheet, View} from "react-native";
import {Select} from "@/components/Select/Select";
import {useEffect} from "react";
import { AppText } from "@/components/ui/AppText";


export interface Props {
    onChange: (value: string) => void
    selected?: string
}

export default function LeagueDropdown(props: Props) {
    const {organization} = useOrganizationStore();
    const {fetchData, data} = useLeagueData()
    const { onChange, selected: selectedProp } = props;

    useEffect(() => {
        if (organization !== null) {
            fetchData({organization_id: organization, limit: 0, offset: 0, term: ""})
        }
    }, [fetchData, organization])

    const onSelectChange = (value: string) => {
        onChange(value)
    }

    const selected = selectedProp ?? data[0]?.id ?? null;

    useEffect(() => {
        if (!selectedProp && data.length > 0) {
            onChange(data[0].id);
        }
    }, [data, onChange, selectedProp]);

    return (
        <View style={[styles.select]}>
            <AppText variant="labelLarge">League</AppText>
            <Select
                onChange={onSelectChange}
                data={data.map(league => ({label: league.name, value: league.id}))}
                selected={selected}
            />
        </View>
    )
}


const styles = StyleSheet.create({
    select: {
        flex: 1,
        gap: 8,
        marginTop: 18
    }
});
