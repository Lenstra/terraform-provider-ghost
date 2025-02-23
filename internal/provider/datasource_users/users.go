package datasource_users

import (
	"context"

	"github.com/Lenstra/terraform-provider-ghost/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func (u *usersDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	users, err := u.client.Users().List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to get users", err.Error())
		return
	}

	state, diags := convert(ctx, users)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func convert(ctx context.Context, users []client.User) (*UsersModel, diag.Diagnostics) {
	elementType := UsersValue{}.Type(ctx)
	if users == nil {
		return &UsersModel{
			Users: types.ListNull(elementType),
		}, nil
	}

	values := make([]UsersValue, len(users))
	for i, user := range users {
		values[i] = UsersValue{
			CommentNotifications:                 basetypes.NewBoolValue(user.CommentNotifications),
			CreatedAt:                            basetypes.NewStringValue(user.CreatedAt),
			DonationNotifications:                basetypes.NewBoolValue(user.DonationNotifications),
			Email:                                basetypes.NewStringValue(user.Email),
			FreeMemberSignupNotification:         basetypes.NewBoolValue(user.FreeMemberSignupNotification),
			Id:                                   basetypes.NewStringValue(user.Id),
			LastSeen:                             basetypes.NewStringValue(user.LastSeen),
			MentionNotifications:                 basetypes.NewBoolValue(user.MentionNotifications),
			MilestoneNotifications:               basetypes.NewBoolValue(user.MilestoneNotifications),
			Name:                                 basetypes.NewStringValue(user.Name),
			PaidSubscriptionCanceledNotification: basetypes.NewBoolValue(user.PaidSubscriptionCanceledNotification),
			PaidSubscriptionStartedNotification:  basetypes.NewBoolValue(user.PaidSubscriptionStartedNotification),
			RecommendationNotifications:          basetypes.NewBoolValue(user.RecommendationNotifications),
			Slug:                                 basetypes.NewStringValue(user.Slug),
			Status:                               basetypes.NewStringValue(user.Status),
			UpdatedAt:                            basetypes.NewStringValue(user.UpdatedAt),
			Url:                                  basetypes.NewStringValue(user.Url),
		}
	}

	value, diags := types.ListValueFrom(ctx, elementType, values)
	if diags.HasError() {
		return nil, diags
	}
	return &UsersModel{Users: value}, diags
}
